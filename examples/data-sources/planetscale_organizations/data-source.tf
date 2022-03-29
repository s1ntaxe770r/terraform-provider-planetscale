data "planetscale_organizations" "orgs" {}

output "organizations" {
  value = data.planetscale_organizations.orgs.organizations
}
