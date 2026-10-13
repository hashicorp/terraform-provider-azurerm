// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_authorization_server

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/authorizationserver"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementAuthorizationServerCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.AuthorizationServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := authorizationserver.NewAuthorizationServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_authorization_server", id.ID())
			}
		}
	}

	authorizationEndpoint := d.Get("authorization_endpoint").(string)
	clientId := d.Get("client_id").(string)
	clientRegistrationEndpoint := d.Get("client_registration_endpoint").(string)
	displayName := d.Get("display_name").(string)
	grantTypesRaw := d.Get("grant_types").(*pluginsdk.Set).List()
	grantTypes := expandApiManagementAuthorizationServerGrantTypes(grantTypesRaw)

	clientAuthenticationMethodsRaw := d.Get("client_authentication_method").(*pluginsdk.Set).List()
	clientAuthenticationMethods := expandApiManagementAuthorizationServerClientAuthenticationMethods(clientAuthenticationMethodsRaw)
	clientSecret := d.Get("client_secret").(string)
	defaultScope := d.Get("default_scope").(string)
	description := d.Get("description").(string)
	resourceOwnerPassword := d.Get("resource_owner_password").(string)
	resourceOwnerUsername := d.Get("resource_owner_username").(string)
	supportState := d.Get("support_state").(bool)
	tokenBodyParametersRaw := d.Get("token_body_parameter").([]any)
	tokenBodyParameters := expandApiManagementAuthorizationServerTokenBodyParameters(tokenBodyParametersRaw)

	params := authorizationserver.AuthorizationServerContract{
		Properties: &authorizationserver.AuthorizationServerContractProperties{
			// Required
			AuthorizationEndpoint:      authorizationEndpoint,
			ClientId:                   clientId,
			ClientRegistrationEndpoint: clientRegistrationEndpoint,
			DisplayName:                displayName,
			GrantTypes:                 pointer.From(grantTypes),

			// Optional
			ClientAuthenticationMethod: clientAuthenticationMethods,
			ClientSecret:               pointer.To(clientSecret),
			DefaultScope:               pointer.To(defaultScope),
			Description:                pointer.To(description),
			ResourceOwnerPassword:      pointer.To(resourceOwnerPassword),
			ResourceOwnerUsername:      pointer.To(resourceOwnerUsername),
			SupportState:               pointer.To(supportState),
			TokenBodyParameters:        tokenBodyParameters,
		},
	}

	authorizationMethodsRaw := d.Get("authorization_methods").(*pluginsdk.Set).List()
	if len(authorizationMethodsRaw) > 0 {
		params.Properties.AuthorizationMethods = expandApiManagementAuthorizationServerAuthorizationMethods(authorizationMethodsRaw)
	}

	bearerTokenSendingMethodsRaw := d.Get("bearer_token_sending_methods").(*pluginsdk.Set).List()
	if len(bearerTokenSendingMethodsRaw) > 0 {
		params.Properties.BearerTokenSendingMethods = expandApiManagementAuthorizationServerBearerTokenSendingMethods(bearerTokenSendingMethodsRaw)
	}

	if tokenEndpoint := d.Get("token_endpoint").(string); tokenEndpoint != "" {
		params.Properties.TokenEndpoint = pointer.To(tokenEndpoint)
	}

	if _, err := client.CreateOrUpdate(ctx, id, params, authorizationserver.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementAuthorizationServerRead(d, meta)
}
