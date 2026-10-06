## v0.1.0 [2026-10-06]

_What's new?_

- New tables added
  - [jenkins_user](docs/tables/jenkins_user.md)

_Enhancements_

- Aligned `jenkins_user` with other Jenkins tables: snake_case columns, struct streaming, logging, `RowsRemaining`, and `id` / `last_change` / `project` columns.
- Removed committed `oryxBuildBinary` build artifact and ignored it going forward.

_Dependencies_

- Recompiled plugin with Go `1.26`.
- Upgraded [steampipe-plugin-sdk](https://github.com/turbot/steampipe-plugin-sdk) from `v5.5.0` to `v5.14.0`.
- Upgraded [gojenkins](https://github.com/IvenGe/gojenkins) from `v1.1.4` to `v1.1.5`.

## v0.0.1 [2023-06-27]

_What's new?_

- New tables added

  - [jenkins_build](https://hub.steampipe.io/plugins/turbot/jenkins/tables/jenkins_build)
  - [jenkins_folder](https://hub.steampipe.io/plugins/turbot/jenkins/tables/jenkins_folder)
  - [jenkins_freestyle_project](https://hub.steampipe.io/plugins/turbot/jenkins/tables/jenkins_freestyle_project)
  - [jenkins_node](https://hub.steampipe.io/plugins/turbot/jenkins/tables/jenkins_node)
  - [jenkins_pipeline](https://hub.steampipe.io/plugins/turbot/jenkins/tables/jenkins_pipeline)
  - [jenkins_plugin](https://hub.steampipe.io/plugins/turbot/jenkins/tables/jenkins_plugin)
