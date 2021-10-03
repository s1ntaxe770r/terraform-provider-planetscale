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

output "all_databases" {
  value = data.planetscale_databases.all.databases  
}




