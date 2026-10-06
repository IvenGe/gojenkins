package jenkins

import (
	"context"
	"strings"
	"time"

	"github.com/IvenGe/gojenkins"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableJenkinsUser() *plugin.Table {
	return &plugin.Table{
		Name:        "jenkins_user",
		Description: "A user account configured in the Jenkins instance.",
		Get: &plugin.GetConfig{
			Hydrate:    getJenkinsUser,
			KeyColumns: plugin.SingleColumn("id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"user not found", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listJenkinsUsers,
		},

		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("User.ID"), Description: "Unique user identifier."},
			{Name: "full_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("User.FullName"), Description: "User's full name."},
			{Name: "absolute_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("User.AbsoluteURL"), Description: "Absolute URL to the user's profile."},
			{Name: "last_change", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LastChange").Transform(epochMsToTime), Description: "Timestamp of the user's last change, when available."},
			{Name: "project", Type: proto.ColumnType_JSON, Description: "Project associated with the user's last change, when available."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("User.FullName"), Description: titleDescription},
		},
	}
}

//// LIST FUNCTION

func listJenkinsUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("jenkins_user.listJenkinsUsers", "connect_error", err)
		return nil, err
	}

	users, err := client.GetAllUsers(ctx)
	if err != nil {
		logger.Error("jenkins_user.listJenkinsUsers", "query_error", err)
		if strings.Contains(err.Error(), "Not found") {
			return nil, nil
		}
		return nil, err
	}

	for _, user := range users.Raw.Users {
		d.StreamListItem(ctx, user)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION

func getJenkinsUser(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Debug("jenkins_user.getJenkinsUser")
	userID := d.EqualsQualString("id")

	// Empty check for userID
	if userID == "" {
		return nil, nil
	}

	client, err := Connect(ctx, d)
	if err != nil {
		logger.Error("jenkins_user.getJenkinsUser", "connect_error", err)
		return nil, err
	}

	user, err := client.GetUser(ctx, userID)
	if err != nil {
		logger.Error("jenkins_user.getJenkinsUser", "query_error", err)
		return nil, err
	}

	// Normalize to the same shape as list rows (PeopleEntry)
	return gojenkins.PeopleEntry{
		User: gojenkins.PeopleUser{
			ID:          user.Raw.ID,
			FullName:    user.Raw.FullName,
			AbsoluteURL: user.Raw.AbsoluteURL,
		},
	}, nil
}

func epochMsToTime(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	var ms int64
	switch v := d.Value.(type) {
	case int64:
		ms = v
	case int:
		ms = int64(v)
	case float64:
		ms = int64(v)
	default:
		return nil, nil
	}

	if ms <= 0 {
		return nil, nil
	}

	return time.UnixMilli(ms).UTC(), nil
}
