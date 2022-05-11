package planetscale

import (
	"context"
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ps "github.com/planetscale/planetscale-go/planetscale"
)

func resourceBranch() *schema.Resource {

	return &schema.Resource{
		Description: "A branch of the database.",

		Create: resourceBranchCreate,
		Read:   resourceBranchRead,
		Delete: resourceBranchDelete,
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Description: "The organization in which the resource belongs.",
				Required:    true,
				ForceNew:    true,
			},
			"database": {
				Type:        schema.TypeString,
				Description: "The database that the branch belongs to.",
				Required:    true,
				ForceNew:    true,
			},
			"parent_branch": {
				Type:        schema.TypeString,
				Description: "The parent branch that the branch is branched from. Default is main.",
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the branch.",
				Required:    true,
				ForceNew:    true,
			},
			"backup_id": {
				Type:        schema.TypeString,
				Description: "The ID of the backup that the branch is branched from.",
				Optional:    true,
				ForceNew:    true,
			},
			"branch": {
				Type:        schema.TypeList,
				Description: "The branch.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the branch.",
							Computed:    true,
						},
						"parent_branch": {
							Type:        schema.TypeString,
							Description: "The parent branch that the branch is branched from. Default is main.",
							Computed:    true,
						},
						"region": {
							Type:        schema.TypeMap,
							Description: "The region the branch is in.",
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ready": {
							Type:        schema.TypeBool,
							Description: "Whether the branch is ready.",
							Computed:    true,
						},
						"production": {
							Type:        schema.TypeBool,
							Description: "Whether the branch is a production branch. ",
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
						"access_host_url": {
							Description: "The host used to connect to the database and branch.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceBranchCreate(d *schema.ResourceData, m interface{}) error {
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
	db, ok := d.GetOk("database")
	if !ok || (db.(string) == "") {
		return errors.New("required value database not set")
	}

	parentBranch := d.Get("parent_branch")

	branch, err := client.DatabaseBranches.Create(ctx, &ps.CreateDatabaseBranchRequest{
		Organization: organization.(string),
		Database:     db.(string),
		Name:         name.(string),
		ParentBranch: parentBranch.(string),
	})

	if err != nil {
		return errors.New(err.Error())
	}

	if err := d.Set("branch", flattenBranch(branch)); err != nil {
		return errors.New(err.Error())
	}
	if err != nil {
		return errors.New("unable to create branch " + err.Error())
	}
	d.SetId(branch.Name)
	return nil
}

func resourceBranchRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ps.Client)
	ctx := context.Background()
	organization, ok := d.GetOk("organization")
	if !ok || (organization.(string) == "") {
		return errors.New("required value organization not set")
	}
	branch, ok := d.GetOk("name")
	if !ok || (branch.(string) == "") {
		return errors.New("required value name not set")
	}
	db, ok := d.GetOk("database")
	if !ok || (db.(string) == "") {
		return errors.New("required value database name not set")
	}

	branchReq, err := client.DatabaseBranches.Get(ctx, &ps.GetDatabaseBranchRequest{
		Organization: organization.(string),
		Database:     db.(string),
		Branch:       "",
	})
	if err != nil {
		// Note(turkenh): Handle "Not Found" gracefully
		if psErr, ok := err.(*ps.Error); ok && psErr.Code == ps.ErrNotFound {
			d.SetId("")
			return nil
		}
		return errors.New(err.Error())
	}
	if err := d.Set("branch", flattenBranch(branchReq)); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func resourceBranchDelete(d *schema.ResourceData, m interface{}) error {
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
	branch, ok := d.GetOk("name")
	if !ok || (branch.(string) == "") {
		return errors.New("required value branch not set")
	}
	err := client.DatabaseBranches.Delete(ctx, &ps.DeleteDatabaseBranchRequest{
		Organization: organization.(string),
		Database:     db.(string),
		Branch:       branch.(string),
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func flattenBranch(branch *ps.DatabaseBranch) (value []map[string]interface{}) {
	if branch != nil {
		v := map[string]interface{}{
			"name":          branch.Name,
			"parent_branch": branch.ParentBranch,
			"region": map[string]interface{}{
				"name":     branch.Region.Name,
				"slug":     branch.Region.Slug,
				"enabled":  strconv.FormatBool(branch.Region.Enabled),
				"location": branch.Region.Location,
			},
			"ready":           branch.Ready,
			"production":      branch.Production,
			"created_at":      branch.CreatedAt.String(),
			"updated_at":      branch.UpdatedAt.String(),
			"access_host_url": branch.AccessHostURL,
		}
		value = append(value, v)
	}
	return value
}
