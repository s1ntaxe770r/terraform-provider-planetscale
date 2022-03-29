data "planetscale_databases" "all" {
  organization = "yourplanescaleorg"
}

output "database" {
  value = data.planetscale_database.all.databases
}
