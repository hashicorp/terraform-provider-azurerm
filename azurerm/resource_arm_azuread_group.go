package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmActiveDirectoryGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmActiveDirectoryGroupCreate,
		Read:   resourceArmActiveDirectoryGroupRead,
		Delete: resourceArmActiveDirectoryGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceArmActiveDirectoryGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).groupsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)

	properties := graphrbac.GroupCreateParameters{
		DisplayName:     &name,
		MailEnabled:     utils.Bool(false),
		MailNickname:    &name,
		SecurityEnabled: utils.Bool(true),
	}

	group, err := client.Create(ctx, properties)
	if err != nil {
		return err
	}

	d.SetId(*group.ObjectID)
	d.Set("object_id", group.ObjectID)

	return resourceArmActiveDirectoryGroupRead(d, meta)
}

func resourceArmActiveDirectoryGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).groupsClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, d.Id())
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] [resource_arm_azuread_group] Azure AD group with id %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Azure AD Group with ID %q: %+v", d.Id(), err)
	}

	d.Set("name", resp.DisplayName)

	return nil
}

func resourceArmActiveDirectoryGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).groupsClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Delete(ctx, d.Id())
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error Deleting Azure AD Group with ID %q: %+v", d.Id(), err)
		}
	}

	return nil
}
