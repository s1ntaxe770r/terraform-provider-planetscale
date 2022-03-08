package planetscale

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ps "github.com/planetscale/planetscale-go/planetscale"
	"strings"
)

func resourceBranchPassword() *schema.Resource {

	return &schema.Resource{
		Create: resourceBranchPasswordCreate,
		Read:   resourceBranchPasswordRead,
		Delete: resourceBranchPasswordDelete,
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Description: "organization to create database under",
				Required:    true,
				ForceNew:    true,
			},
			"database": {
				Type:        schema.TypeString,
				Description: "display name of your database",
				Required:    true,
				ForceNew:    true,
			},
			"branch": {
				Type:        schema.TypeString,
				Description: "display name of your branch",
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "name of your password",
				Required:    true,
				ForceNew:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "id of your password",
				Computed:    true,
				ForceNew:    false,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "creation time of your password",
				Computed:    true,
				ForceNew:    false,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "your plain password",
				Sensitive:   true,
				Computed:    true,
				ForceNew:    false,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "username to connect database through host",
				Sensitive:   false,
				Computed:    true,
				ForceNew:    false,
			},
			"host": {
				Type:        schema.TypeString,
				Description: "host to connect database branch",
				Sensitive:   false,
				Computed:    true,
				ForceNew:    false,
			},
		},
	}
}

func resourceBranchPasswordCreate(d *schema.ResourceData, m interface{}) error {
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
	branch, ok := d.GetOk("branch")
	if !ok || (branch.(string) == "") {
		return errors.New("required value branch not set")
	}
	name, ok := d.GetOk("name")
	if !ok || (name.(string) == "") {
		return errors.New("required value name not set")
	}

	password, err := client.Passwords.Create(ctx, &ps.DatabaseBranchPasswordRequest{
		Organization: organization.(string),
		Database:     db.(string),
		Branch:       branch.(string),
		DisplayName:  name.(string),
	})

	if err != nil {
		return errors.New(err.Error())
	}

	err = d.Set("id", password.PublicID)
	if err != nil {
		return errors.New("unable to create password " + err.Error())
	}
	err = d.Set("password", password.PlainText)
	if err != nil {
		return errors.New("unable to create password " + err.Error())
	}
	err = d.Set("created_at", password.CreatedAt.String())
	if err != nil {
		return errors.New("unable to create password " + err.Error())
	}

	connectionString := password.ConnectionStrings.General

	if !ok || (connectionString == "") {
		return errors.New("api returned empty connection string. Connection string can not be empty")
	}

	elementMap := createElementMapFromConnectionStrings(connectionString)

	if len(elementMap) == 0 {
		return errors.New("connection string can not be empty")
	}

	err = d.Set("username", elementMap["username"])
	if err != nil {
		return errors.New("username can not be empty " + err.Error())
	}

	err = d.Set("host", elementMap["host"])
	if err != nil {
		return errors.New("host can not be empty" + err.Error())
	}

	d.SetId(password.PublicID)

	return nil
}

func resourceBranchPasswordRead(d *schema.ResourceData, m interface{}) error {
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
	branch, ok := d.GetOk("branch")
	if !ok || (branch.(string) == "") {
		return errors.New("required value branch not set")
	}
	name, ok := d.GetOk("name")
	if !ok || (name.(string) == "") {
		return errors.New("required value name not set")
	}
	id, ok := d.GetOk("id")
	if !ok || (id.(string) == "") {
		return errors.New("required value name not set")
	}

	password, err := client.Passwords.Get(ctx, &ps.GetDatabaseBranchPasswordRequest{
		Organization: organization.(string),
		Database:     db.(string),
		Branch:       branch.(string),
		DisplayName:  name.(string),
		PasswordId:   id.(string),
	})
	if err != nil {
		return errors.New(err.Error())
	}

	d.SetId(password.PublicID)

	err = d.Set("id", password.PublicID)
	if err != nil {
		return errors.New("unable to create password " + err.Error())
	}

	// passwords can not be displayed after creation
	// I will remain this block to prevent others to add a similar statement
	// err = d.Set("password", password.PlainText)

	err = d.Set("created_at", password.CreatedAt.String())
	if err != nil {
		return errors.New("unable to create password " + err.Error())
	}

	connectionString := password.ConnectionStrings.General
	if !ok || (connectionString == "") {
		return errors.New("connection String can not be empty")
	}

	elementMap := createElementMapFromConnectionStrings(connectionString)

	if len(elementMap) == 0 {
		return errors.New("connection string can not be empty")
	}

	err = d.Set("username", elementMap["username"])
	if err != nil {
		return errors.New("username can not be empty " + err.Error())
	}

	err = d.Set("host", elementMap["host"])
	if err != nil {
		return errors.New("host can not be empty" + err.Error())
	}

	return nil
}

func createElementMapFromConnectionStrings(s string) map[string]string {
	elementMap := make(map[string]string)
	for _, data := range strings.Split(s, "\n") {
		if len(data) > 0 {
			tokens := strings.Split(data, ":")
			elementMap[strings.TrimSpace(tokens[0])] = strings.TrimSpace(tokens[1])
		}
	}
	return elementMap
}

func resourceBranchPasswordDelete(d *schema.ResourceData, m interface{}) error {
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
	branch, ok := d.GetOk("branch")
	if !ok || (branch.(string) == "") {
		return errors.New("required value branch not set")
	}
	name, ok := d.GetOk("name")
	if !ok || (name.(string) == "") {
		return errors.New("required value name not set")
	}
	id, ok := d.GetOk("id")
	if !ok || (id.(string) == "") {
		return errors.New("required value ID not set")
	}
	err := client.Passwords.Delete(ctx, &ps.DeleteDatabaseBranchPasswordRequest{
		Organization: organization.(string),
		Database:     db.(string),
		Branch:       branch.(string),
		DisplayName:  name.(string),
		PasswordId:   id.(string),
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
