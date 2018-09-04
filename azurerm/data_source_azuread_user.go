package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
)

func dataSourceArmAzureADUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAzureADUserRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"user_principal_name"},
			},

			"user_principal_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"object_id"},
			},
		},
	}
}

func dataSourceArmAzureADUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).usersClient
	ctx := meta.(*ArmClient).StopContext

	var user graphrbac.User
	var queryString string

	if oId, ok := d.GetOk("object_id"); ok {

		// use the object_id to find the Azure AD user
		queryString = oId.(string)
	} else {

		// use the user_principal_name to find the Azure AD user
		queryString = d.Get("user_principal_name").(string)
	}

	log.Printf("[DEBUG] [data_source_azuread_user] Using Get with upnOrObjectId: %q", queryString)
	resp, err := client.Get(ctx, queryString)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: AzureAD User with ID or userPrincipalName %q not found", queryString)
		}
		return fmt.Errorf("Error making Read request on AzureAD User with ID or userPrincipalName %q: %+v", queryString, err)
	}

	user = resp

	d.SetId(*user.ObjectID)
	d.Set("object_id", user.ObjectID)
	d.Set("user_principal_name", user.UserPrincipalName)

	return nil
}
