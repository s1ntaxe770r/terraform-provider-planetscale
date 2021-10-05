package planetscale

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ps "github.com/planetscale/planetscale-go/planetscale"
)

func dataSourceOrganizations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOrganizationsRead,
		Schema: map[string]*schema.Schema{
			"organizations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOrganizationsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ps.Client)
	ctx := context.Background()
	organizationslist, err := client.Organizations.List(ctx)
	if err != nil {
		return errors.New("unable to list organizations " + err.Error())
	}
	organizations := make([]map[string]interface{}, 0)
	for _, organization := range flattenOrganizations(organizationslist) {
		organizations = append(organizations, organization)
	}
	if err := d.Set("organizations", organizations); err != nil {
		return err
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}

func flattenOrganizations(organizations []*ps.Organization) (values []map[string]interface{}) {
	if organizations != nil {
		for _, organization := range organizations {
			v := map[string]interface{}{
				"name":       organization.Name,
				"created_at": organization.CreatedAt.String(),
				"updated_at": organization.UpdatedAt.String(),
			}
			values = append(values, v)
		}
	}
	return values
}
