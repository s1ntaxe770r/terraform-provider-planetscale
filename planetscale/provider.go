package planetscale

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/planetscale/planetscale-go/planetscale"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PLANETSCALE_TOKEN", nil),
			},
		},
		ConfigureFunc: configureProvider,
		ResourcesMap:  map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"planetscale_databases": dataSourceDatabases(),
			"planetscale_database":  dataSourceDatabase(),
		},
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	access_token := d.Get("access_token")
	c, err := planetscale.NewClient(planetscale.WithAccessToken(access_token.(string)))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to create planetscale client, %s", err.Error()))
	}
	return c, nil
}
