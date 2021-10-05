package planetscale

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ps "github.com/planetscale/planetscale-go/planetscale"
)

func dataSourceDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDatabaseRead,
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db": {
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"slug": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"display_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"location": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enabled": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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

func dataSourceDatabaseRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ps.Client)
	ctx := context.Background()
	organization, ok := d.GetOk("organization")
	if !ok || (organization.(string) == "") {
		return errors.New("required value organization not set")
	}
	db, ok := d.GetOk("database")
	if !ok || (db.(string) == "") {
		return errors.New("required value database not set")
	}

	databaseresp, err := client.Databases.Get(ctx, &ps.GetDatabaseRequest{
		Organization: organization.(string),
		Database:     db.(string),
	})
	if err != nil {
		return errors.New(err.Error())
	}

	if err := d.Set("db", flattenDatabase(databaseresp)); err != nil {
		return errors.New(err.Error())
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}

func flattenDatabase(database *ps.Database) (value []map[string]interface{}) {
	if database != nil {
		v := map[string]interface{}{
			"name":  database.Name,
			"notes": database.Notes,
			"region": map[string]interface{}{
				"name":     database.Region.Name,
				"slug":     database.Region.Slug,
				"enabled":  strconv.FormatBool(database.Region.Enabled),
				"location": database.Region.Location,
			},
			"created_at": database.CreatedAt.String(),
			"updated_at": database.UpdatedAt.String(),
		}
		value = append(value, v)
	}
	return value
}
