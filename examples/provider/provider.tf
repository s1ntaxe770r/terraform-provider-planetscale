provider "planetscale" {
  access_token = var.access_token # optionally use PLANETSCALE_ACCESS_TOKEN env var

  service_token    = var.service_token    # optionally use PLANETSCALE_SERVICE_TOKEN env var
  service_token_id = var.service_token_id # optionally use PLANETSCALE_SERVICE_TOKEN_ID env var
}
