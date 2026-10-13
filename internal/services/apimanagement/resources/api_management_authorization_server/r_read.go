// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_authorization_server

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/authorizationserver"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementAuthorizationServerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.AuthorizationServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := authorizationserver.ParseAuthorizationServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("api_management_name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("name", pointer.From(model.Name))
		if props := model.Properties; props != nil {
			d.Set("authorization_endpoint", props.AuthorizationEndpoint)
			d.Set("client_id", props.ClientId)
			d.Set("client_registration_endpoint", props.ClientRegistrationEndpoint)
			d.Set("default_scope", pointer.From(props.DefaultScope))
			d.Set("description", pointer.From(props.Description))
			d.Set("display_name", props.DisplayName)
			d.Set("support_state", pointer.From(props.SupportState))
			d.Set("token_endpoint", pointer.From(props.TokenEndpoint))

			// TODO: Read properties from api, https://github.com/Azure/azure-rest-api-specs/issues/14128
			d.Set("resource_owner_password", d.Get("resource_owner_password").(string))
			d.Set("resource_owner_username", d.Get("resource_owner_username").(string))

			if err := d.Set("authorization_methods", flattenApiManagementAuthorizationServerAuthorizationMethods(props.AuthorizationMethods)); err != nil {
				return fmt.Errorf("flattening `authorization_methods`: %+v", err)
			}

			if err := d.Set("bearer_token_sending_methods", flattenApiManagementAuthorizationServerBearerTokenSendingMethods(props.BearerTokenSendingMethods)); err != nil {
				return fmt.Errorf("flattening `bearer_token_sending_methods`: %+v", err)
			}

			if err := d.Set("client_authentication_method", flattenApiManagementAuthorizationServerClientAuthenticationMethods(props.ClientAuthenticationMethod)); err != nil {
				return fmt.Errorf("flattening `client_authentication_method`: %+v", err)
			}

			if err := d.Set("grant_types", flattenApiManagementAuthorizationServerGrantTypes(props.GrantTypes)); err != nil {
				return fmt.Errorf("flattening `grant_types`: %+v", err)
			}

			if err := d.Set("token_body_parameter", flattenApiManagementAuthorizationServerTokenBodyParameters(props.TokenBodyParameters)); err != nil {
				return fmt.Errorf("flattening `token_body_parameter`: %+v", err)
			}
		}
	}

	return nil
}
