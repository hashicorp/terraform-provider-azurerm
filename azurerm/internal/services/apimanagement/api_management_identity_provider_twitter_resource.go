package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementIdentityProviderTwitter() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiManagementIdentityProviderTwitterCreateUpdate,
		Read:   resourceApiManagementIdentityProviderTwitterRead,
		Update: resourceApiManagementIdentityProviderTwitterCreateUpdate,
		Delete: resourceApiManagementIdentityProviderTwitterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"api_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"api_secret_key": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceApiManagementIdentityProviderTwitterCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	clientID := d.Get("api_key").(string)
	clientSecret := d.Get("api_secret_key").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, apimanagement.Twitter)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Identity Provider %q (API Management Service %q / Resource Group %q): %s", apimanagement.Twitter, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_identity_provider_twitter", *existing.ID)
		}
	}

	parameters := apimanagement.IdentityProviderCreateContract{
		IdentityProviderCreateContractProperties: &apimanagement.IdentityProviderCreateContractProperties{
			ClientID:     utils.String(clientID),
			ClientSecret: utils.String(clientSecret),
			Type:         apimanagement.Twitter,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apimanagement.Twitter, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating Identity Provider %q (Resource Group %q / API Management Service %q): %+v", apimanagement.Twitter, resourceGroup, serviceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, apimanagement.Twitter)
	if err != nil {
		return fmt.Errorf("retrieving Identity Provider %q (Resource Group %q / API Management Service %q): %+v", apimanagement.Twitter, resourceGroup, serviceName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Identity Provider %q (Resource Group %q / API Management Service %q)", apimanagement.Twitter, resourceGroup, serviceName)
	}
	d.SetId(*resp.ID)

	return resourceApiManagementIdentityProviderTwitterRead(d, meta)
}

func resourceApiManagementIdentityProviderTwitterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IdentityProviderID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	identityProviderName := id.Name

	resp, err := client.Get(ctx, resourceGroup, serviceName, apimanagement.IdentityProviderType(identityProviderName))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Identity Provider %q (Resource Group %q / API Management Service %q) was not found - removing from state!", identityProviderName, resourceGroup, serviceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for Identity Provider %q (Resource Group %q / API Management Service %q): %+v", identityProviderName, resourceGroup, serviceName, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	if props := resp.IdentityProviderContractProperties; props != nil {
		d.Set("api_key", props.ClientID)
	}

	return nil
}

func resourceApiManagementIdentityProviderTwitterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IdentityProviderID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	identityProviderName := id.Name

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, apimanagement.IdentityProviderType(identityProviderName), ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Identity Provider %q (Resource Group %q / API Management Service %q): %+v", identityProviderName, resourceGroup, serviceName, err)
		}
	}

	return nil
}
