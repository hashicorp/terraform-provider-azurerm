// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	path := d.Get("path").(string)
	version := d.Get("version").(string)
	versionSetId := d.Get("version_set_id").(string)
	displayName := d.Get("display_name").(string)
	protocolsRaw := d.Get("protocols").(*pluginsdk.Set).List()
	protocols := expandApiManagementApiProtocols(protocolsRaw)
	sourceApiId := d.Get("source_api_id").(string)
	serviceUrl := d.Get("service_url").(string)

	id, err := api.ParseApiID(d.Id())
	if err != nil {
		return err
	}

	apiType := api.ApiTypeHTTP
	if v, ok := d.GetOk("api_type"); ok {
		apiType = api.ApiType(v.(string))
	}
	soapApiType := soapApiTypeFromApiType(apiType)

	// If import is used, we need to send properties to Azure API in two operations.
	// First we execute import and then updated the other props.
	if d.HasChange("import") {
		if vs, hasImport := d.GetOk("import"); hasImport {
			if apiParams := expandApiManagementApiImport(vs.([]any), apiType, soapApiType,
				path, serviceUrl, version, versionSetId); apiParams != nil {
				result, err := client.CreateOrUpdate(ctx, *id, *apiParams, api.CreateOrUpdateOperationOptions{})
				if err != nil {
					return fmt.Errorf("creating with import of %s: %+v", id, err)
				}

				if pollerType := custompollers.NewAPIManagementAPIPoller(client, *id, result.HttpResponse); pollerType != nil {
					poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
					if err := poller.PollUntilDone(ctx); err != nil {
						return fmt.Errorf("polling import %s: %+v", id, err)
					}
				}
			}
		}
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	existing := resp.Model.Properties
	if existing.Type != nil {
		soapApiType = soapApiTypeFromApiType(pointer.From(existing.Type))
	}
	prop := &api.ApiCreateOrUpdateProperties{
		Path:                          existing.Path,
		Protocols:                     existing.Protocols,
		ServiceURL:                    existing.ServiceURL,
		Description:                   existing.Description,
		ApiVersionDescription:         existing.ApiVersionDescription,
		ApiRevisionDescription:        existing.ApiRevisionDescription,
		SubscriptionRequired:          existing.SubscriptionRequired,
		SubscriptionKeyParameterNames: existing.SubscriptionKeyParameterNames,
		Contact:                       existing.Contact,
		License:                       existing.License,
		SourceApiId:                   existing.SourceApiId,
		DisplayName:                   existing.DisplayName,
		ApiVersion:                    existing.ApiVersion,
		ApiVersionSetId:               existing.ApiVersionSetId,
		TermsOfServiceURL:             existing.TermsOfServiceURL,
		Type:                          existing.Type,
		ApiType:                       pointer.To(soapApiType),
	}

	// For the setting of `AuthenticationSettingsContract`, the PUT payload restrictions are as follows:
	//   1. Cannot have both 'oAuth2' and 'openid' set
	//   2. Cannot use `OAuth2AuthenticationSettings` in combination with `OAuth2` nor `openid`
	//   3. Cannot use `OpenidAuthenticationSettings` in combination with `Openid` nor `OAuth2`
	// If specifying `oauth2_authorization`/`openid_authentication` when creating a resource and then updating the resource, the error #2/#3 mentioned above will occur.
	// This is because starting from the 2022-08-01 version, the Get API additionally returns a collection of `oauth2_authorization`/`openid_authentication` authentication settings, which property name is `OAuth2AuthenticationSettings`/`OpenidAuthenticationSetting`.
	// Given the API behaviour, the update here should only read the specified property `oauth2_authorization`/`openid_authentication` to exclude `OAuth2AuthenticationSettings`/`OpenidAuthenticationSetting` to ensure the update works properly.
	if v := existing.AuthenticationSettings; v != nil {
		authenticationSettings := &api.AuthenticationSettingsContract{}
		if v.OAuth2 != nil {
			authenticationSettings.OAuth2 = v.OAuth2
			prop.AuthenticationSettings = authenticationSettings
		}

		if v.Openid != nil {
			authenticationSettings.Openid = v.Openid
			prop.AuthenticationSettings = authenticationSettings
		}
	}

	if d.HasChange("path") {
		prop.Path = path
	}

	if d.HasChange("protocols") {
		prop.Protocols = protocols
	}

	if d.HasChange("api_type") {
		prop.Type = pointer.To(apiType)
		prop.ApiType = pointer.To(soapApiType)
	}

	if d.HasChange("service_url") {
		prop.ServiceURL = pointer.To(serviceUrl)
	}

	if d.HasChange("description") {
		prop.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("revision_description") {
		prop.ApiRevisionDescription = pointer.To(d.Get("revision_description").(string))
	}

	if d.HasChange("version_description") {
		prop.ApiVersionDescription = pointer.To(d.Get("version_description").(string))
	}
	if d.HasChange("subscription_required") {
		prop.SubscriptionRequired = pointer.To(d.Get("subscription_required").(bool))
	}

	if d.HasChange("subscription_key_parameter_names") {
		subscriptionKeyParameterNamesRaw := d.Get("subscription_key_parameter_names").([]any)
		prop.SubscriptionKeyParameterNames = expandApiManagementApiSubscriptionKeyParamNames(subscriptionKeyParameterNamesRaw)
	}

	if d.HasChange("oauth2_authorization") {
		authenticationSettings := &api.AuthenticationSettingsContract{}
		oAuth2AuthorizationSettingsRaw := d.Get("oauth2_authorization").([]any)
		authenticationSettings.OAuth2 = expandApiManagementOAuth2AuthenticationSettingsContract(oAuth2AuthorizationSettingsRaw)
		prop.AuthenticationSettings = authenticationSettings
	}

	if d.HasChange("openid_authentication") {
		authenticationSettings := &api.AuthenticationSettingsContract{}
		openIDAuthorizationSettingsRaw := d.Get("openid_authentication").([]any)
		authenticationSettings.Openid = expandApiManagementOpenIDAuthenticationSettingsContract(openIDAuthorizationSettingsRaw)
		prop.AuthenticationSettings = authenticationSettings
	}

	if d.HasChange("contact") {
		prop.Contact = expandApiManagementApiContact(d.Get("contact").([]any))
	}

	if d.HasChange("license") {
		prop.License = expandApiManagementApiLicense(d.Get("license").([]any))
	}

	if d.HasChange("source_api_id") {
		prop.SourceApiId = pointer.To(sourceApiId)
	}

	if d.HasChange("display_name") {
		prop.DisplayName = pointer.To(displayName)
	}

	if d.HasChange("version") {
		prop.ApiVersion = pointer.To(version)
	}

	if d.HasChange("version_set_id") {
		prop.ApiVersionSetId = pointer.To(versionSetId)
	}

	if d.HasChange("terms_of_service_url") {
		prop.TermsOfServiceURL = pointer.To(d.Get("terms_of_service_url").(string))
	}

	params := api.ApiCreateOrUpdateParameter{
		Properties: prop,
	}

	result, err := client.CreateOrUpdate(ctx, *id, params, api.CreateOrUpdateOperationOptions{IfMatch: pointer.To("*")})
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if pollerType := custompollers.NewAPIManagementAPIPoller(client, *id, result.HttpResponse); pollerType != nil {
		poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("polling creating/updating %s: %+v", id, err)
		}
	}

	return resourceApiManagementApiRead(d, meta)
}
