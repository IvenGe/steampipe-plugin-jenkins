# Table: jenkins_user

A Jenkins user account known to the controller.

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

### Find a user by name

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
