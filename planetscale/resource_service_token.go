package planetscale

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ps "github.com/planetscale/planetscale-go/planetscale"
)

func resourceServiceToken() *schema.Resource {
	return &schema.Resource{
		Create: ServiceTokenCreate,
		Delete: ServiceTokenDelete,
		Read:   ServiceTokenRead,
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"token": {
				Type:        schema.TypeString,
				Description: "token returned upon creation. Pls note this a sensitive",
				Computed:    true,
				Sensitive:   true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "id of the token created",
				Computed:    true,
			},
		},
	}
}

func ServiceTokenCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ps.Client)
	ctx := context.Background()
	organization, ok := d.GetOk("organization")
	if !ok || (organization.(string) == "") {
		return errors.New("required value organization not set")
	}
	token, err := client.ServiceTokens.Create(ctx, &ps.CreateServiceTokenRequest{
		Organization: organization.(string),
	})
	err = d.Set("token", token.Token)
	if err != nil {
		return err
	}
	err = d.Set("id", token.ID)
	if err != nil {
		return err
	}
	return nil
}

func ServiceTokenDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*ps.Client)
	ctx := context.Background()
	organization, ok := d.GetOk("organization")
	if !ok || (organization.(string) == "") {
		return errors.New("required value organization not set")
	}
	id := d.Get("id")
	err := client.ServiceTokens.Delete(ctx, &ps.DeleteServiceTokenRequest{
		Organization: organization.(string),
		ID:           id.(string),
	})
	if err != nil {
		return err
	}
	return nil
}

func ServiceTokenRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ps.Client)
	ctx := context.Background()
	organization, ok := d.GetOk("organization")
	if !ok || (organization.(string) == "") {
		return errors.New("required value organization not set")
	}
	id := d.Get("id")

	tokens, err := client.ServiceTokens.List(ctx, &ps.ListServiceTokensRequest{
		Organization: organization.(string),
	})
	if err != nil {
		return errors.New(fmt.Sprintf("unable to list service tokens, %s", err.Error()))
	}
	for _, v := range tokens {
		if v.ID == id.(string) {
			if err := d.Set("token", v.Token); err != nil {
				return err
			}
			if err := d.Set("id", v.ID); err != nil {
				return err
			}
		}
	}
	return nil
}
