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
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/authorizationserver"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementAuthorizationServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementAuthorizationServerCreateUpdate,
		Read:   resourceApiManagementAuthorizationServerRead,
		Update: resourceApiManagementAuthorizationServerCreateUpdate,
		Delete: resourceApiManagementAuthorizationServerDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := authorizationserver.ParseAuthorizationServerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"authorization_endpoint": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"authorization_methods": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(authorizationserver.AuthorizationMethodDELETE),
						string(authorizationserver.AuthorizationMethodGET),
						string(authorizationserver.AuthorizationMethodHEAD),
						string(authorizationserver.AuthorizationMethodOPTIONS),
						string(authorizationserver.AuthorizationMethodPATCH),
						string(authorizationserver.AuthorizationMethodPOST),
						string(authorizationserver.AuthorizationMethodPUT),
						string(authorizationserver.AuthorizationMethodTRACE),
					}, false),
				},
				Set: pluginsdk.HashString,
			},

			"client_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"client_registration_endpoint": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"grant_types": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(authorizationserver.GrantTypeAuthorizationCode),
						string(authorizationserver.GrantTypeClientCredentials),
						string(authorizationserver.GrantTypeImplicit),
						string(authorizationserver.GrantTypeResourceOwnerPassword),
					}, false),
				},
				Set: pluginsdk.HashString,
			},

			// Optional
			"bearer_token_sending_methods": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(authorizationserver.BearerTokenSendingMethodAuthorizationHeader),
						string(authorizationserver.BearerTokenSendingMethodQuery),
					}, false),
				},
				Set: pluginsdk.HashString,
			},

			"client_authentication_method": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(authorizationserver.ClientAuthenticationMethodBasic),
						string(authorizationserver.ClientAuthenticationMethodBody),
					}, false),
				},
				Set: pluginsdk.HashString,
			},

			"client_secret": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"default_scope": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"resource_owner_username": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"resource_owner_password": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"support_state": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"token_body_parameter": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"value": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"token_endpoint": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceApiManagementAuthorizationServerCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.AuthorizationServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := authorizationserver.NewAuthorizationServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
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
	tokenBodyParametersRaw := d.Get("token_body_parameter").([]interface{})
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
		authorizationMethods := expandApiManagementAuthorizationServerAuthorizationMethods(authorizationMethodsRaw)
		params.Properties.AuthorizationMethods = authorizationMethods
	}

	bearerTokenSendingMethodsRaw := d.Get("bearer_token_sending_methods").(*pluginsdk.Set).List()
	if len(bearerTokenSendingMethodsRaw) > 0 {
		bearerTokenSendingMethods := expandApiManagementAuthorizationServerBearerTokenSendingMethods(bearerTokenSendingMethodsRaw)
		params.Properties.BearerTokenSendingMethods = bearerTokenSendingMethods
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

func resourceApiManagementAuthorizationServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceApiManagementAuthorizationServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.AuthorizationServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := authorizationserver.ParseAuthorizationServerID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, authorizationserver.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %s", *id, err)
		}
	}

	return nil
}

func expandApiManagementAuthorizationServerGrantTypes(input []interface{}) *[]authorizationserver.GrantType {
	outputs := make([]authorizationserver.GrantType, 0)

	for _, v := range input {
		grantType := authorizationserver.GrantType(v.(string))
		outputs = append(outputs, grantType)
	}

	return &outputs
}

func flattenApiManagementAuthorizationServerGrantTypes(input []authorizationserver.GrantType) []interface{} {
	outputs := make([]interface{}, 0)

	for _, v := range input {
		outputs = append(outputs, string(v))
	}

	return outputs
}

func expandApiManagementAuthorizationServerAuthorizationMethods(input []interface{}) *[]authorizationserver.AuthorizationMethod {
	outputs := make([]authorizationserver.AuthorizationMethod, 0)

	for _, v := range input {
		grantType := authorizationserver.AuthorizationMethod(v.(string))
		outputs = append(outputs, grantType)
	}

	return &outputs
}

func flattenApiManagementAuthorizationServerAuthorizationMethods(input *[]authorizationserver.AuthorizationMethod) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	for _, v := range *input {
		outputs = append(outputs, string(v))
	}

	return outputs
}

func expandApiManagementAuthorizationServerBearerTokenSendingMethods(input []interface{}) *[]authorizationserver.BearerTokenSendingMethod {
	outputs := make([]authorizationserver.BearerTokenSendingMethod, 0)

	for _, v := range input {
		grantType := authorizationserver.BearerTokenSendingMethod(v.(string))
		outputs = append(outputs, grantType)
	}

	return &outputs
}

func flattenApiManagementAuthorizationServerBearerTokenSendingMethods(input *[]authorizationserver.BearerTokenSendingMethod) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	for _, v := range *input {
		outputs = append(outputs, string(v))
	}

	return outputs
}

func expandApiManagementAuthorizationServerClientAuthenticationMethods(input []interface{}) *[]authorizationserver.ClientAuthenticationMethod {
	outputs := make([]authorizationserver.ClientAuthenticationMethod, 0)

	for _, v := range input {
		grantType := authorizationserver.ClientAuthenticationMethod(v.(string))
		outputs = append(outputs, grantType)
	}

	return &outputs
}

func flattenApiManagementAuthorizationServerClientAuthenticationMethods(input *[]authorizationserver.ClientAuthenticationMethod) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	for _, v := range *input {
		outputs = append(outputs, string(v))
	}

	return outputs
}

func expandApiManagementAuthorizationServerTokenBodyParameters(input []interface{}) *[]authorizationserver.TokenBodyParameterContract {
	outputs := make([]authorizationserver.TokenBodyParameterContract, 0)

	for _, v := range input {
		vs := v.(map[string]interface{})

		output := authorizationserver.TokenBodyParameterContract{
			Name:  vs["name"].(string),
			Value: vs["value"].(string),
		}
		outputs = append(outputs, output)
	}

	return &outputs
}

func flattenApiManagementAuthorizationServerTokenBodyParameters(input *[]authorizationserver.TokenBodyParameterContract) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	for _, v := range *input {
		output := make(map[string]interface{})

		output["name"] = v.Name
		output["value"] = v.Value

		outputs = append(outputs, output)
	}

	return outputs
}
