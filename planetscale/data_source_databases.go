package planetscale

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ps "github.com/planetscale/planetscale-go/planetscale"
)

func dataSourceDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDatabasesRead,
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notes": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "represents the state of a database",
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

func dataSourceDatabasesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ps.Client)
	ctx := context.Background()
	organization, ok := d.GetOk("organization")
	if !ok || (organization.(string) == "") {
		return errors.New("required value organization not set")
	}

	databaselist, err := client.Databases.List(ctx, &ps.ListDatabasesRequest{
		Organization: organization.(string),
	})
	if err != nil {
		return err
	}
	databases := make([]map[string]interface{}, 0)
	for _, db := range flattenDatabases(databaselist) {
		databases = append(databases, db)
	}
	if err := d.Set("databases", databases); err != nil {
		return err
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}

func flattenDatabases(databases []*ps.Database) (values []map[string]interface{}) {
	if databases != nil {
		for _, database := range databases {
			v := map[string]interface{}{
				"name":  database.Name,
				"notes": database.Notes,
				"region": map[string]interface{}{
					"name":     database.Region.Name,
					"slug":     database.Region.Slug,
					"enabled":  strconv.FormatBool(database.Region.Enabled),
					"location": database.Region.Location,
				},
				"state":      string(database.State),
				"created_at": database.CreatedAt.String(),
				"updated_at": database.UpdatedAt.String(),
			}
			values = append(values, v)
		}
	}
	return values
}
