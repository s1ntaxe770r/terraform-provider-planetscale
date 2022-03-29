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
				Type:          schema.TypeString,
				Description:   "User access token. Can be specified with the `PLANETSCALE_ACCESS_TOKEN` environment variable.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("PLANETSCALE_ACCESS_TOKEN", nil),
				ConflictsWith: []string{"service_token", "service_token_id"},
			},
			"service_token": {
				Type:          schema.TypeString,
				Description:   "Service token. Can be specified with the `PLANETSCALE_SERVICE_TOKEN` environment variable.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("PLANETSCALE_SERVICE_TOKEN", nil),
				ConflictsWith: []string{"access_token"},
			},
			"service_token_id": {
				Type:          schema.TypeString,
				Description:   "ID generated alongside a service token. Can be specified with the `PLANETSCALE_SERVICE_TOKEN_ID` environment variable.",
				Optional:      true,
				Sensitive:     true,
				DefaultFunc:   schema.EnvDefaultFunc("PLANETSCALE_SERVICE_TOKEN_ID", nil),
				ConflictsWith: []string{"access_token"},
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
	accessToken := d.Get("access_token").(string)
	opt := planetscale.WithAccessToken(accessToken)

	serviceToken, tokenOk := d.GetOk("service_token")
	serviceTokenID, idOk := d.GetOk("service_token_id")

	if tokenOk || idOk {
		if !tokenOk {
			return nil, diag.FromErr(fmt.Errorf("service_token must be set if service_token_id is set"))
		}
		if !idOk {
			return nil, diag.FromErr(fmt.Errorf("service_token_id must be set if service_token is set"))
		}

		opt = planetscale.WithServiceToken(serviceTokenID.(string), serviceToken.(string))
	}

	c, err := planetscale.NewClient(opt)
	if err != nil {
		return nil, diag.FromErr(errors.New(fmt.Sprintf("unable to create planetscale client, %s", err.Error())))
	}
	return c, nil
}
