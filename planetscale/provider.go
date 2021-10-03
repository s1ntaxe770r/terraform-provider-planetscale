package planetscale

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/planetscale/planetscale-go/planetscale"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PLANETSCALE_TOKEN", nil),
			},
		},
		ConfigureFunc: configureProvider,
		ResourcesMap:  map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"planetscale_databases": dataSourceDatabases(),
		},
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	access_token := d.Get("access_token")
	c, err := planetscale.NewClient(planetscale.WithAccessToken(access_token.(string)))
	if err != nil {
		log.Fatalf("Unable to create planetscale client  %s", err.Error())
		return nil, err
	}
	return c, nil
}
