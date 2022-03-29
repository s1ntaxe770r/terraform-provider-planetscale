resource "planetscale_branch_password" "instance" {
  organization = var.db.org
  database     = var.db.name
  branch       = var.db.branch
  name         = "${terraform.workspace}-tf"
}
