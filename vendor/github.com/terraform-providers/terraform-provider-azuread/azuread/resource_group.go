package azuread

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/graph"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/p"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Delete: resourceGroupDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).groupsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)

	properties := graphrbac.GroupCreateParameters{
		DisplayName:     &name,
		MailEnabled:     p.Bool(false),                 // we're defaulting to false, as the API currently only supports the creation of non-mail enabled security groups.
		MailNickname:    p.String(uuid.New().String()), // this matches the portal behaviour
		SecurityEnabled: p.Bool(true),                  // we're defaulting to true, as the API currently only supports the creation of non-mail enabled security groups.
	}

	group, err := client.Create(ctx, properties)
	if err != nil {
		return fmt.Errorf("Error creating Group (%q): %+v", name, err)
	}
	if group.ObjectID == nil {
		return fmt.Errorf("nil Group ID for %q: %+v", name, err)
	}
	d.SetId(*group.ObjectID)

	_, err = graph.WaitForReplication(func() (interface{}, error) {
		return client.Get(ctx, *group.ObjectID)
	})
	if err != nil {
		return fmt.Errorf("Error waiting for Group (%s) with ObjectId %q: %+v", name, *group.ObjectID, err)
	}

	return resourceGroupRead(d, meta)
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).groupsClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, d.Id())
	if err != nil {
		if ar.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Azure AD group with id %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Azure AD Group with ID %q: %+v", d.Id(), err)
	}

	d.Set("name", resp.DisplayName)
	d.Set("object_id", resp.ObjectID)
	return nil
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).groupsClient
	ctx := meta.(*ArmClient).StopContext

	if resp, err := client.Delete(ctx, d.Id()); err != nil {
		if !ar.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error Deleting Azure AD Group with ID %q: %+v", d.Id(), err)
		}
	}

	return nil
}
