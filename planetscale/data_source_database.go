package planetscale

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ps "github.com/planetscale/planetscale-go/planetscale"
)

func dataSourceDatabase() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve a database by name.",

		Read: dataSourceDatabaseRead,
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Description: "The organization in which the resource belongs.",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the database.",
				Required:    true,
			},
			"notes": {
				Type:        schema.TypeString,
				Description: "A description of the database.",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeMap,
				Description: "The region the database.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"state": {
				Type:        schema.TypeString,
				Description: "The state of the database instance.",
				Computed:    true,
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
	}
}

func dataSourceDatabaseRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ps.Client)
	ctx := context.Background()
	organization, ok := d.GetOk("organization")
	if !ok || (organization.(string) == "") {
		return errors.New("required value organization not set")
	}
	db, ok := d.GetOk("name")
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
	if err := d.Set("region", flattenRegion(&databaseresp.Region)); err != nil {
		return errors.New(err.Error())
	}
	if err := d.Set("notes", databaseresp.Notes); err != nil {
		return errors.New(err.Error())
	}
	if err := d.Set("created_at", databaseresp.CreatedAt.String()); err != nil {
		return errors.New(err.Error())
	}
	if err := d.Set("updated_at", databaseresp.UpdatedAt.String()); err != nil {
		return errors.New(err.Error())
	}
	if err := d.Set("state", string(databaseresp.State)); err != nil {
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
			"state":      string(database.State),
			"created_at": database.CreatedAt.String(),
			"updated_at": database.UpdatedAt.String(),
		}
		value = append(value, v)
	}
	return value
}

func flattenRegion(region *ps.Region) (value map[string]string) {
	if region != nil {
		v := map[string]string{
			"name":     region.Name,
			"slug":     region.Slug,
			"enabled":  strconv.FormatBool(region.Enabled),
			"location": region.Location,
		}
		value = v
	}
	return value
}
