package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmAzureADGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAzureADGroupRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"object_id"},
			},

			"object_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validate.UUID,
			},
		},
	}
}

func dataSourceArmAzureADGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).groupsClient
	ctx := meta.(*ArmClient).StopContext

	var adgroup graphrbac.ADGroup
	var groupObj *graphrbac.ADGroup

	if oId, ok := d.GetOk("object_id"); ok {
		// use the object_id to find the Azure AD group

		objectId := oId.(string)
		resp, err := client.Get(ctx, objectId)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error: AzureAD Group with ID %q was not found", objectId)
			}

			return fmt.Errorf("Error making Read request on AzureAD Group with ID %q: %+v", objectId, err)
		}

		adgroup = resp

	} else {

		// use the name to find the Azure AD group
		name := d.Get("name").(string)
		filter := fmt.Sprintf("displayName eq '%s'", name)
		log.Printf("[DEBUG] [data_source_azuread_group] Using filter %q", filter)

		resp, err := client.ListComplete(ctx, filter)
		if err != nil {
			return fmt.Errorf("Error listing Azure AD groups: %+v", err)
		}

		for _, v := range *resp.Response().Value {
			if v.DisplayName != nil {
				if strings.EqualFold(*v.DisplayName, name) {
					log.Printf("[DEBUG] [data_source_azuread_group] %q (API result) matches %q (given value). The group has the objectId: %q", *v.DisplayName, name, *v.ObjectID)
					groupObj = &v
					break
				} else {
					log.Printf("[DEBUG] [data_source_azuread_group] %q (API result) does not match %q (given value)", *v.DisplayName, name)
				}
			}
		}
		if groupObj == nil {
			return fmt.Errorf("Couldn't locate a Azure AD group with a name of %q", name)
		}

		adgroup = *groupObj
	}

	d.SetId(*adgroup.ObjectID)
	d.Set("object_id", adgroup.ObjectID)
	d.Set("name", adgroup.DisplayName)

	return nil
}
