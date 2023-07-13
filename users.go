package gojenkins

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	createUserContext = "/securityRealm/createAccountByAdmin"
)

// User is a Jenkins account
type User struct {
	Jenkins  *Jenkins
	UserName string
	FullName string
	Email    string
	Raw      *UserRespone
}

type Users struct {
	Jenkins  *Jenkins
	UserName string
	FullName string
	Email    string
	Id       string
	Base     string
	Raw      *UserRespone
}

type UserRespone struct {
	//
	Class       string `json:"_class"`
	AbsoluteUrl string `json:"absoluteUrl"`
	Description string `json:"description"`
	FullName    string `json:"fullName"`
	Id          string `json:"id"`
}

// ErrUser occurs when there is error creating or revoking Jenkins users
type ErrUser struct {
	Message string
}

func (e *ErrUser) Error() string {
	return e.Message
}

// CreateUser creates a new Jenkins account
func (j *Jenkins) CreateUser(ctx context.Context, userName, password, fullName, email string) (User, error) {
	user := User{
		// Set Jenkins client pointer to be able to delete user later
		Jenkins:  j,
		UserName: userName,
		FullName: fullName,
		Email:    email,
	}
	payload := "username=" + userName + "&password1=" + password + "&password2=" + password + "&fullname=" + fullName + "&email=" + email
	response, err := j.Requester.Post(ctx, createUserContext, strings.NewReader(payload), nil, nil)
	if err != nil {
		return user, err
	}
	if response.StatusCode != http.StatusOK {
		return user, &ErrUser{
			Message: fmt.Sprintf("error creating user. Status is %d", response.StatusCode),
		}
	}
	return user, nil
}

// DeleteUser deletes a Jenkins account
func (j *Jenkins) DeleteUser(ctx context.Context, userName string) error {
	deleteContext := "/securityRealm/user/" + userName + "/doDelete"
	payload := "Submit=Yes"
	response, err := j.Requester.Post(ctx, deleteContext, strings.NewReader(payload), nil, nil)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return &ErrUser{
			Message: fmt.Sprintf("error deleting user. Status is %d", response.StatusCode),
		}
	}
	return nil
}

// Delete deletes a Jenkins account
func (u *User) Delete() error {
	return u.Jenkins.DeleteUser(context.Background(), u.UserName)
}

var userAPIResp struct {
	UserName string `json:"id"`
	FullName string `json:"fullName"`
}

func (u *User) GetUser(ctx context.Context, userName string) (User, error) {
	getUserContext := "/securityRealm/user/" + userName + "/api/json"
	response, err := u.Jenkins.Requester.GetJSON(ctx, getUserContext, nil, nil)
	if err != nil {
		return User{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return User{}, errors.New(fmt.Sprintf("error retrieving user. Status is %d", response.StatusCode))
	}

	err = json.NewDecoder(response.Body).Decode(&userAPIResp)
	if err != nil {
		return User{}, err
	}

	user := User{
		Jenkins:  u.Jenkins,
		UserName: userAPIResp.UserName,
		FullName: userAPIResp.FullName,
	}

	return user, nil
}

func (u *Users) Poll(ctx context.Context) (int, error) {

	response, err := u.Jenkins.Requester.GetJSON(ctx, u.Base, u.Raw, nil)
	if err != nil {
		return 0, err
	}
	return response.StatusCode, nil
}
