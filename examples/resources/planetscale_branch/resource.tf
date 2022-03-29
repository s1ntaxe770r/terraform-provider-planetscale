resource "planetscale_branch" "db1" {
  organization  = "startupheroes"
  database      = "my-test-database"
  parent_branch = "main"
  name          = "my-pretty-fetaure"
}
