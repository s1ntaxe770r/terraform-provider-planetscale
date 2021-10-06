terraform {
  required_providers {
    planetscale = {
      version = "0.3.1"
      source = "jubril.me/edu/planetscale"
    }
  }
}

provider "planetscale" {
  access_token = "pscale_oauth_o6xxDnY-IfaRnp-55GJQf2tbKr7DS1xhpfBaCeMKnOA"
}
# data "planetscale_databases" "all" {
#   organization = "gophercorp"
# }
data "planetscale_database" "one"{
  organization = "gophercorp"
  name = "express"
}

# data "planetscale_organizations" "orgs" {
  
# }

# output "all_databases" {
#   value = data.planetscale_databases.all.databases  
# }

# output "database" {
#   value = data.planetscale_database.one.db
# }
# output "organizations" {
#   value = data.planetscale_organizations.orgs.organizations
# }

resource "planetscale_database" "primarydb"{
  organization = "gophercorp"
  name = "tftest"
}
