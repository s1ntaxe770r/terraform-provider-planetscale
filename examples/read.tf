terraform {
  required_providers {
    planetscale = {
      version = "0.3.1"
      source = "jubril.me/edu/planetscale"
    }
  }
}

provider "planetscale" {
  access_token = ""
}

data "planetscale_databases" "all" {
  organization = "gophercorp"
}

data "planetscale_database" "one"{
  organization = "gophercorp"
  database = "express"
}
output "all_databases" {
  value = data.planetscale_databases.all.databases  
}

output "database" {
  value = data.planetscale_database.one.db
}




