package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementIdentityProviderAADB2C() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementIdentityProviderAADB2CCreateUpdate,
		Read:   resourceArmApiManagementIdentityProviderAADB2CRead,
		Update: resourceArmApiManagementIdentityProviderAADB2CCreateUpdate,
		Delete: resourceArmApiManagementIdentityProviderAADB2CDelete,
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

			"client_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"client_secret": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// For AADB2C identity providers, `allowed_tenants` must specify exactly one tenant
			"allowed_tenant": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"signin_tenant": {
				Type:     schema.TypeString,
				Required: true,
				// B2C tenant domains can be customized, and GUIDs might work here too
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"authority": {
				Type:     schema.TypeString,
				Required: true,
				// B2C login domains can be customized and don't necessarily end in b2clogin.com
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"signup_policy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"signin_policy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"profile_editing_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"password_reset_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceArmApiManagementIdentityProviderAADB2CCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)

	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	allowedTenant := d.Get("allowed_tenant").(string)
	signinTenant := d.Get("signin_tenant").(string)
	authority := d.Get("authority").(string)
	signupPolicy := d.Get("signup_policy").(string)

	signinPolicy := d.Get("signin_policy").(string)
	profileEditingPolicy := d.Get("profile_editing_policy").(string)
	passwordResetPolicy := d.Get("password_reset_policy").(string)

	id := parse.NewIdentityProviderID(client.SubscriptionID, resourceGroup, serviceName, string(apimanagement.AadB2C))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, apimanagement.AadB2C)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id.String(), err)
			}
		} else {
			return tf.ImportAsExistsError("azurerm_api_management_identity_provider_aadb2c", id.ID())
		}
	}

	parameters := apimanagement.IdentityProviderCreateContract{
		IdentityProviderCreateContractProperties: &apimanagement.IdentityProviderCreateContractProperties{
			ClientID:                 utils.String(clientID),
			ClientSecret:             utils.String(clientSecret),
			Type:                     apimanagement.AadB2C,
			AllowedTenants:           utils.ExpandStringSlice([]interface{}{allowedTenant}),
			SigninTenant:             utils.String(signinTenant),
			Authority:                utils.String(authority),
			SignupPolicyName:         utils.String(signupPolicy),
			SigninPolicyName:         utils.String(signinPolicy),
			ProfileEditingPolicyName: utils.String(profileEditingPolicy),
			PasswordResetPolicyName:  utils.String(passwordResetPolicy),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apimanagement.AadB2C, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating Identity Provider %q (Resource Group %q / API Management Service %q): %+v", apimanagement.AadB2C, resourceGroup, serviceName, err)
	}

	d.SetId(id.ID())
	return resourceArmApiManagementIdentityProviderAADB2CRead(d, meta)
}

func resourceArmApiManagementIdentityProviderAADB2CRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IdentityProviderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, apimanagement.IdentityProviderType(id.Name))
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Identity Provider %q (Resource Group %q / API Management Service %q) was not found - removing from state!", id.Name, id.ResourceGroup, id.ServiceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for Identity Provider %q (Resource Group %q / API Management Service %q): %+v", id.Name, id.ResourceGroup, id.ServiceName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("api_management_name", id.ServiceName)

	if props := resp.IdentityProviderContractProperties; props != nil {
		d.Set("client_id", props.ClientID)
		d.Set("signin_tenant", props.SigninTenant)
		d.Set("authority", props.Authority)
		d.Set("signup_policy", props.SignupPolicyName)
		d.Set("signin_policy", props.SigninPolicyName)
		d.Set("profile_editing_policy", props.ProfileEditingPolicyName)
		d.Set("password_reset_policy", props.PasswordResetPolicyName)

		allowedTenant := ""
		if allowedTenants := props.AllowedTenants; allowedTenants != nil && len(*allowedTenants) > 0 {
			t := *allowedTenants
			allowedTenant = t[0]
		}
		d.Set("allowed_tenant", allowedTenant)
	}

	return nil
}

func resourceArmApiManagementIdentityProviderAADB2CDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IdentityProviderID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, apimanagement.IdentityProviderType(id.Name), ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Identity Provider %q (Resource Group %q / API Management Service %q): %+v", id.Name, id.ResourceGroup, id.ServiceName, err)
		}
	}

	return nil
}
