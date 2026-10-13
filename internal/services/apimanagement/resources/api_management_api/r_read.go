// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := api.ParseApiID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("api_management_name", id.ServiceName)
	d.Set("name", helpers.GetApiName(id.ApiId))
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			apiType := pointer.FromEnum(props.Type)
			if len(apiType) == 0 {
				apiType = string(api.ApiTypeHTTP)
			}
			d.Set("api_type", apiType)
			d.Set("description", pointer.From(props.Description))
			d.Set("display_name", pointer.From(props.DisplayName))
			d.Set("is_current", pointer.From(props.IsCurrent))
			d.Set("is_online", pointer.From(props.IsOnline))
			d.Set("path", props.Path)
			d.Set("service_url", pointer.From(props.ServiceURL))
			d.Set("revision", pointer.From(props.ApiRevision))
			d.Set("subscription_required", pointer.From(props.SubscriptionRequired))
			d.Set("version", pointer.From(props.ApiVersion))
			d.Set("version_set_id", pointer.From(props.ApiVersionSetId))
			d.Set("revision_description", pointer.From(props.ApiRevisionDescription))
			d.Set("version_description", pointer.From(props.ApiVersionDescription))
			d.Set("terms_of_service_url", pointer.From(props.TermsOfServiceURL))

			if err := d.Set("protocols", flattenApiManagementApiProtocols(props.Protocols)); err != nil {
				return fmt.Errorf("setting `protocols`: %s", err)
			}

			if err := d.Set("subscription_key_parameter_names", flattenApiManagementApiSubscriptionKeyParamNames(props.SubscriptionKeyParameterNames)); err != nil {
				return fmt.Errorf("setting `subscription_key_parameter_names`: %+v", err)
			}

			if err := d.Set("oauth2_authorization", flattenApiManagementOAuth2Authorization(props.AuthenticationSettings.OAuth2)); err != nil {
				return fmt.Errorf("setting `oauth2_authorization`: %+v", err)
			}

			if err := d.Set("openid_authentication", flattenApiManagementOpenIDAuthentication(props.AuthenticationSettings.Openid)); err != nil {
				return fmt.Errorf("setting `openid_authentication`: %+v", err)
			}

			if err := d.Set("contact", flattenApiManagementApiContact(props.Contact)); err != nil {
				return fmt.Errorf("setting `contact`: %+v", err)
			}

			if err := d.Set("license", flattenApiManagementApiLicense(props.License)); err != nil {
				return fmt.Errorf("setting `license`: %+v", err)
			}
		}
	}
	return nil
}

func soapApiTypeFromApiType(apiType api.ApiType) api.SoapApiType {
	return map[api.ApiType]api.SoapApiType{
		api.ApiTypeGraphql:   api.SoapApiTypeGraphql,
		api.ApiTypeHTTP:      api.SoapApiTypeHTTP,
		api.ApiTypeSoap:      api.SoapApiTypeSoap,
		api.ApiTypeWebsocket: api.SoapApiTypeWebsocket,
	}[apiType]
}
