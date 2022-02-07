package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementIdentityProviderAAD() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementIdentityProviderAADCreateUpdate,
		Read:   resourceApiManagementIdentityProviderAADRead,
		Update: resourceApiManagementIdentityProviderAADCreateUpdate,
		Delete: resourceApiManagementIdentityProviderAADDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"client_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"client_secret": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"allowed_tenants": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},
			"signin_tenant": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceApiManagementIdentityProviderAADCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	allowedTenants := d.Get("allowed_tenants").([]interface{})
	signinTenant := d.Get("signin_tenant").(string)
	id := parse.NewIdentityProviderID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), string(apimanagement.IdentityProviderTypeAad))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, apimanagement.IdentityProviderTypeAad)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_identity_provider_aad", id.ID())
		}
	}

	parameters := apimanagement.IdentityProviderCreateContract{
		IdentityProviderCreateContractProperties: &apimanagement.IdentityProviderCreateContractProperties{
			ClientID:       utils.String(clientID),
			ClientSecret:   utils.String(clientSecret),
			Type:           apimanagement.IdentityProviderTypeAad,
			AllowedTenants: utils.ExpandStringSlice(allowedTenants),
			SigninTenant:   utils.String(signinTenant),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, apimanagement.IdentityProviderTypeAad, parameters, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementIdentityProviderAADRead(d, meta)
}

func resourceApiManagementIdentityProviderAADRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		d.Set("client_id", props.ClientID)
		d.Set("allowed_tenants", props.AllowedTenants)
		d.Set("signin_tenant", props.SigninTenant)
	}

	return nil
}

func resourceApiManagementIdentityProviderAADDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
