package planetscale

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/planetscale/planetscale-go/planetscale"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("PLANETSCALE_ACCESS_TOKEN", nil),
			},
		},
		ConfigureContextFunc: configureProvider,
		ResourcesMap: map[string]*schema.Resource{
			"planetscale_database":        resourceDatabase(),
			"planetscale_branch":          resourceBranch(),
			"planetscale_branch_password": resourceBranchPassword(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"planetscale_databases":     dataSourceDatabases(),
			"planetscale_database":      dataSourceDatabase(),
			"planetscale_organizations": dataSourceOrganizations(),
		},
	}
}

func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	accessToken := d.Get("access_token")
	c, err := planetscale.NewClient(planetscale.WithAccessToken(accessToken.(string)))
	if err != nil {
		return nil, diag.FromErr(errors.New(fmt.Sprintf("unable to create planetscale client, %s", err.Error())))
	}
	return c, nil
}
