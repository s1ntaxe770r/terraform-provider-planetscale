resource "planetscale_database" "db" {
  organization = "exampleorg"
  name         = "exampledb"
}

output "database" {
  value = data.planetscale_database.db.region
}
