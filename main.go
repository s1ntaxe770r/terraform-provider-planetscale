package main

import (
	"terraform-provider-planetscale/planetscale"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: planetscale.Provider,
	})
}
