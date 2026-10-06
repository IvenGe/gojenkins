# Table: jenkins_user

A user account configured in the Jenkins instance.

## Examples

### List all users

```sql
select
  id,
  full_name,
  absolute_url
from
  jenkins_user
order by
  full_name;
```

### Find users by name

```sql
select
  id,
  full_name,
  absolute_url
from
  jenkins_user
where
  full_name ilike '%admin%';
```

### Users with recent changes

```sql
select
  id,
  full_name,
  last_change,
  project ->> 'name' as project_name
from
  jenkins_user
where
  last_change is not null
order by
  last_change desc;
```

### Get a specific user by id

```sql
select
  id,
  full_name,
  absolute_url,
  last_change
from
  jenkins_user
where
  id = 'admin';
```
