// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/identityprovider"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmApiManagementIdentityProviderAADB2C() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmApiManagementIdentityProviderAADB2CCreateUpdate,
		Read:   resourceArmApiManagementIdentityProviderAADB2CRead,
		Update: resourceArmApiManagementIdentityProviderAADB2CCreateUpdate,
		Delete: resourceArmApiManagementIdentityProviderAADB2CDelete,

		Importer: identityProviderImportFunc(identityprovider.IdentityProviderTypeAadBTwoC),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

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

			// For AADB2C identity providers, `allowed_tenants` must specify exactly one tenant
			"allowed_tenant": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"signin_tenant": {
				Type:     pluginsdk.TypeString,
				Required: true,
				// B2C tenant domains can be customized, and GUIDs might work here too
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"authority": {
				Type:     pluginsdk.TypeString,
				Required: true,
				// B2C login domains can be customized and don't necessarily end in b2clogin.com
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"signup_policy": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"signin_policy": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"client_library": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 16),
			},

			"profile_editing_policy": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"password_reset_policy": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceArmApiManagementIdentityProviderAADB2CCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)

	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	clientLibrary := d.Get("client_library").(string)

	allowedTenant := d.Get("allowed_tenant").(string)
	signinTenant := d.Get("signin_tenant").(string)
	authority := d.Get("authority").(string)
	signupPolicy := d.Get("signup_policy").(string)

	signinPolicy := d.Get("signin_policy").(string)
	profileEditingPolicy := d.Get("profile_editing_policy").(string)
	passwordResetPolicy := d.Get("password_reset_policy").(string)

	id := identityprovider.NewIdentityProviderID(meta.(*clients.Client).Account.SubscriptionId, resourceGroup, serviceName, identityprovider.IdentityProviderTypeAadBTwoC)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id.String(), err)
			}
		} else {
			return tf.ImportAsExistsError("azurerm_api_management_identity_provider_aadb2c", id.ID())
		}
	}

	parameters := identityprovider.IdentityProviderCreateContract{
		Properties: &identityprovider.IdentityProviderCreateContractProperties{
			ClientId:                 clientID,
			ClientLibrary:            pointer.To(clientLibrary),
			ClientSecret:             clientSecret,
			Type:                     pointer.To(identityprovider.IdentityProviderTypeAadBTwoC),
			AllowedTenants:           utils.ExpandStringSlice([]interface{}{allowedTenant}),
			SigninTenant:             pointer.To(signinTenant),
			Authority:                pointer.To(authority),
			SignupPolicyName:         pointer.To(signupPolicy),
			SigninPolicyName:         pointer.To(signinPolicy),
			ProfileEditingPolicyName: pointer.To(profileEditingPolicy),
			PasswordResetPolicyName:  pointer.To(passwordResetPolicy),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, identityprovider.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceArmApiManagementIdentityProviderAADB2CRead(d, meta)
}

func resourceArmApiManagementIdentityProviderAADB2CRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := identityprovider.ParseIdentityProviderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("client_id", props.ClientId)
			d.Set("client_library", pointer.From(props.ClientLibrary))
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
	}

	return nil
}

func resourceArmApiManagementIdentityProviderAADB2CDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := identityprovider.ParseIdentityProviderID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, identityprovider.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
