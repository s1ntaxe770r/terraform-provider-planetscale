data "planetscale_database" "db" {
  organization = "yourplanescaleorg"
  name         = "your planetscaledb"
}

output "database" {
  value = data.planetscale_database.db.region
}
