// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_identity_provider_aadb2c

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/identityprovider"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmApiManagementIdentityProviderAADB2CCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id.String(), err)
				}
			} else {
				return tf.ImportAsExistsError("azurerm_api_management_identity_provider_aadb2c", id.ID())
			}
		}
	}

	parameters := identityprovider.IdentityProviderCreateContract{
		Properties: &identityprovider.IdentityProviderCreateContractProperties{
			ClientId:                 clientID,
			ClientLibrary:            pointer.To(clientLibrary),
			ClientSecret:             clientSecret,
			Type:                     pointer.To(identityprovider.IdentityProviderTypeAadBTwoC),
			AllowedTenants:           pluginsdk.ExpandStringSlice([]any{allowedTenant}),
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
