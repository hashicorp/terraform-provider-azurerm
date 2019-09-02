package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmActiveDirectoryServicePrincipal() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: `The Azure Active Directory resources have been split out into their own Provider.

Information on migrating to the new AzureAD Provider can be found here: https://terraform.io/docs/providers/azurerm/guides/migrating-to-azuread.html

As such the Azure Active Directory resources within the AzureRM Provider are now deprecated and will be removed in v2.0 of the AzureRM Provider.
`,
		Read: dataSourceArmActiveDirectoryServicePrincipalRead,

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"display_name", "application_id"},
			},

			"display_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"object_id", "application_id"},
			},

			"application_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"object_id", "display_name"},
			},
		},
	}
}

func dataSourceArmActiveDirectoryServicePrincipalRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).graph.ServicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	var servicePrincipal *graphrbac.ServicePrincipal

	if v, ok := d.GetOk("object_id"); ok {

		//use the object_id to find the Azure AD service principal
		objectId := v.(string)
		app, err := client.Get(ctx, objectId)
		if err != nil {
			if utils.ResponseWasNotFound(app.Response) {
				return fmt.Errorf("Service Principal with Object ID %q was not found!", objectId)
			}

			return fmt.Errorf("Error retrieving Service Principal ID %q: %+v", objectId, err)
		}

		servicePrincipal = &app

	} else if _, ok := d.GetOk("display_name"); ok {

		// use the display_name to find the Azure AD service principal
		displayName := d.Get("display_name").(string)
		filter := fmt.Sprintf("displayName eq '%s'", displayName)
		log.Printf("[DEBUG] [data_source_azuread_service_principal] Using filter %q", filter)

		apps, err := client.ListComplete(ctx, filter)
		if err != nil {
			return fmt.Errorf("Error listing Service Principals: %+v", err)
		}

		for _, app := range *apps.Response().Value {
			if app.DisplayName == nil {
				continue
			}

			if *app.DisplayName == displayName {
				servicePrincipal = &app
				break
			}
		}

		if servicePrincipal == nil {
			return fmt.Errorf("A Service Principal with the Display Name %q was not found", displayName)
		}

	} else {

		// use the application_id to find the Azure AD service principal
		applicationId := d.Get("application_id").(string)
		filter := fmt.Sprintf("appId eq '%s'", applicationId)
		log.Printf("[DEBUG] [data_source_azuread_service_principal] Using filter %q", filter)

		apps, err := client.ListComplete(ctx, filter)
		if err != nil {
			return fmt.Errorf("Error listing Service Principals: %+v", err)
		}

		for _, app := range *apps.Response().Value {
			if app.AppID == nil {
				continue
			}

			if *app.AppID == applicationId {
				servicePrincipal = &app
				break
			}
		}

		if servicePrincipal == nil {
			return fmt.Errorf("A Service Principal for Application ID %q was not found", applicationId)
		}

	}

	d.SetId(*servicePrincipal.ObjectID)

	d.Set("application_id", servicePrincipal.AppID)
	d.Set("display_name", servicePrincipal.DisplayName)
	d.Set("object_id", servicePrincipal.ObjectID)

	return nil
}
