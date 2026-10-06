## v0.1.0 [2026-10-06]

_What's new?_

- Reliability and correctness improvements across all tables
  - Correct Jenkins "not found" error matching (`404`, `Build not found`, `No node found`, …)
  - Propagate folder recursion errors instead of silently dropping nested jobs
  - Timeout-aware HTTP client for Jenkins connections
  - `jenkins_build` Get hydrate now includes `job_full_name`
  - Repaired `jenkins_user` columns (`id`, `full_name`, `absolute_url`), description, and docs
  - Safer nil checks and shared job path splitting helpers
- Dependency: bump `github.com/IvenGe/gojenkins` to **v1.1.6**
- Add unit tests for path splitting, not-found predicates, and config casting

## v0.0.1 [2023-06-27]

_What's new?_

- New tables added

  - [jenkins_build](https://hub.steampipe.io/plugins/turbot/jenkins/tables/jenkins_build)
  - [jenkins_folder](https://hub.steampipe.io/plugins/turbot/jenkins/tables/jenkins_folder)
  - [jenkins_freestyle_project](https://hub.steampipe.io/plugins/turbot/jenkins/tables/jenkins_freestyle_project)
  - [jenkins_node](https://hub.steampipe.io/plugins/turbot/jenkins/tables/jenkins_node)
  - [jenkins_pipeline](https://hub.steampipe.io/plugins/turbot/jenkins/tables/jenkins_pipeline)
  - [jenkins_plugin](https://hub.steampipe.io/plugins/turbot/jenkins/tables/jenkins_plugin)
