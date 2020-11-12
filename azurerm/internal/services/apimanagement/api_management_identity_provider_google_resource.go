package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementIdentityProviderGoogle() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementIdentityProviderGoogleCreateUpdate,
		Read:   resourceArmApiManagementIdentityProviderGoogleRead,
		Update: resourceArmApiManagementIdentityProviderGoogleCreateUpdate,
		Delete: resourceArmApiManagementIdentityProviderGoogleDelete,
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
			"resource_group_name": azure.SchemaResourceGroupNameDeprecated(), // TODO remove in v3.0

			"api_management_name": azure.SchemaApiManagementNameDeprecated(), // TODO remove in v3.0

			"api_management_id": {
				Type:         schema.TypeString,
				Optional:     true, // TODO change to required in v3.0
				Computed:     true, // TODO remove in v3.0
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"client_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.GoogleClientID,
			},

			"client_secret": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceArmApiManagementIdentityProviderGoogleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var resourceGroup, serviceName string
	if apiManagementId, ok := d.GetOk("api_management_id"); ok && apiManagementId.(string) != "" {
		id, err := parse.ApiManagementID(apiManagementId.(string))
		if err != nil {
			return err
		}
		resourceGroup = id.ResourceGroup
		serviceName = id.ServiceName
	} else {
		resourceGroup = d.Get("resource_group_name").(string)
		serviceName = d.Get("api_management_name").(string)
	}

	if resourceGroup == "" || serviceName == "" {
		return fmt.Errorf("could not determine resource group or API management service, please specify `api_management_id`")
	}

	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, apimanagement.Google)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Identity Provider %q (API Management Service %q / Resource Group %q): %s", apimanagement.Google, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_identity_provider_google", *existing.ID)
		}
	}

	parameters := apimanagement.IdentityProviderCreateContract{
		IdentityProviderCreateContractProperties: &apimanagement.IdentityProviderCreateContractProperties{
			ClientID:     utils.String(clientID),
			ClientSecret: utils.String(clientSecret),
			Type:         apimanagement.Google,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apimanagement.Google, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating Identity Provider %q (Resource Group %q / API Management Service %q): %+v", apimanagement.Google, resourceGroup, serviceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, apimanagement.Google)
	if err != nil {
		return fmt.Errorf("retrieving Identity Provider %q (Resource Group %q / API Management Service %q): %+v", apimanagement.Google, resourceGroup, serviceName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Identity Provider %q (Resource Group %q / API Management Service %q)", apimanagement.Google, resourceGroup, serviceName)
	}
	d.SetId(*resp.ID)

	return resourceArmApiManagementIdentityProviderGoogleRead(d, meta)
}

func resourceArmApiManagementIdentityProviderGoogleRead(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiManagementIdentityProviderID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	identityProviderName := id.ProviderName

	resp, err := client.Get(ctx, resourceGroup, serviceName, apimanagement.IdentityProviderType(identityProviderName))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Identity Provider %q (Resource Group %q / API Management Service %q) was not found - removing from state!", identityProviderName, resourceGroup, serviceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for Identity Provider %q (Resource Group %q / API Management Service %q): %+v", identityProviderName, resourceGroup, serviceName, err)
	}

	d.Set("api_management_id", id.ApiManagementID(subscriptionId))

	if props := resp.IdentityProviderContractProperties; props != nil {
		d.Set("client_id", props.ClientID)
	}

	return nil
}

func resourceArmApiManagementIdentityProviderGoogleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	identityProviderName := id.Path["identityProviders"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, apimanagement.IdentityProviderType(identityProviderName), ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Identity Provider %q (Resource Group %q / API Management Service %q): %+v", identityProviderName, resourceGroup, serviceName, err)
		}
	}

	return nil
}
