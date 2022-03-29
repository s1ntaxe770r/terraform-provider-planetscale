package planetscale

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ps "github.com/planetscale/planetscale-go/planetscale"
)

func resourceDatabase() *schema.Resource {
	return &schema.Resource{
		Description: "A PlanetScale database.",

		Create: resourceDatabaseCreate,
		Read:   resourceDatabaseRead,
		Delete: resourceDatabaseDelete,
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Description: "The organization in which the resource belongs.",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the database.",
				Required:    true,
				ForceNew:    true,
			},
			"notes": {
				Type:        schema.TypeString,
				Description: "A description of the database.",
				Optional:    true,
				ForceNew:    true,
			},
			"database": {
				Type:        schema.TypeList,
				Description: "The database.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the database.",
							Computed:    true,
						},
						"notes": {
							Type:        schema.TypeString,
							Description: "A description of the database.",
							Computed:    true,
						},
						"region": {
							Type:        schema.TypeMap,
							Description: "The region the database will sit in.",
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
				},
			},
		},
	}
}

func resourceDatabaseCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ps.Client)
	ctx := context.Background()
	organization, ok := d.GetOk("organization")
	if !ok || (organization.(string) == "") {
		return errors.New("required value organization not set")
	}
	name, ok := d.GetOk("name")
	if !ok || (name.(string) == "") {
		return errors.New("required value name not set")
	}
	notes := d.Get("notes")
	db, err := client.Databases.Create(ctx, &ps.CreateDatabaseRequest{
		Name:         name.(string),
		Organization: organization.(string),
		Notes:        notes.(string),
	})
	if err := d.Set("database", flattenDatabase(db)); err != nil {
		return errors.New(err.Error())
	}
	if err != nil {
		return errors.New("unable to create database " + err.Error())
	}
	d.SetId(db.Name)
	return nil
}

func resourceDatabaseRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ps.Client)
	ctx := context.Background()
	organization, ok := d.GetOk("organization")
	if !ok || (organization.(string) == "") {
		return errors.New("required value organization not set")
	}
	db, ok := d.GetOk("name")
	if !ok || (db.(string) == "") {
		return errors.New("required value database name not set")
	}

	databaseresp, err := client.Databases.Get(ctx, &ps.GetDatabaseRequest{
		Organization: organization.(string),
		Database:     db.(string),
	})
	if err != nil {
		return errors.New(err.Error())
	}
	if err := d.Set("database", flattenDatabase(databaseresp)); err != nil {
		return errors.New(err.Error())
	}
	return nil
}
func resourceDatabaseDelete(d *schema.ResourceData, m interface{}) error {
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
	err := client.Databases.Delete(ctx, &ps.DeleteDatabaseRequest{
		Organization: organization.(string),
		Database:     db.(string),
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
