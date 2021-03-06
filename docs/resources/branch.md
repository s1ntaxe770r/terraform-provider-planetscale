---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "planetscale_branch Resource - terraform-provider-planetscale"
subcategory: ""
description: |-
  A branch of the database.
---

# planetscale_branch (Resource)

A branch of the database.

## Example Usage

```terraform
resource "planetscale_branch" "db1" {
  organization  = "startupheroes"
  database      = "my-test-database"
  parent_branch = "main"
  name          = "my-pretty-fetaure"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `database` (String) The database that the branch belongs to.
- `name` (String) The name of the branch.
- `organization` (String) The organization in which the resource belongs.

### Optional

- `backup_id` (String) The ID of the backup that the branch is branched from.
- `id` (String) The ID of this resource.
- `parent_branch` (String) The parent branch that the branch is branched from. Default is main.

### Read-Only

- `branch` (List of Object) The branch. (see [below for nested schema](#nestedatt--branch))

<a id="nestedatt--branch"></a>
### Nested Schema for `branch`

Read-Only:

- `access_host_url` (String)
- `created_at` (String)
- `name` (String)
- `parent_branch` (String)
- `production` (Boolean)
- `ready` (Boolean)
- `region` (Map of String)
- `updated_at` (String)


