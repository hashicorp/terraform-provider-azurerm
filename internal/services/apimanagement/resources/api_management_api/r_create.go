// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	revision := d.Get("revision").(string)
	path := d.Get("path").(string)
	apiId := fmt.Sprintf("%s;rev=%s", d.Get("name").(string), revision)
	version := d.Get("version").(string)
	versionSetId := d.Get("version_set_id").(string)
	displayName := d.Get("display_name").(string)
	protocolsRaw := d.Get("protocols").(*pluginsdk.Set).List()
	protocols := expandApiManagementApiProtocols(protocolsRaw)
	sourceApiId := d.Get("source_api_id").(string)

	id := api.NewApiID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), apiId)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of an existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_api", id.ID())
		}
	}

	apiType := api.ApiTypeHTTP
	if v, ok := d.GetOk("api_type"); ok {
		apiType = api.ApiType(v.(string))
	}
	soapApiType := soapApiTypeFromApiType(apiType)

	// If import is used, we need to send properties to Azure API in two operations.
	// First we execute import and then updated the other props.
	if importVs, ok := d.GetOk("import"); ok {
		if apiParams := expandApiManagementApiImport(importVs.([]any), apiType, soapApiType,
			path, d.Get("service_url").(string), version, versionSetId); apiParams != nil {
			result, err := client.CreateOrUpdate(ctx, id, *apiParams, api.CreateOrUpdateOperationOptions{})
			if err != nil {
				return fmt.Errorf("creating with import of %s: %+v", id, err)
			}

			if pollerType := custompollers.NewAPIManagementAPIPoller(client, id, result.HttpResponse); pollerType != nil {
				poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("polling import %s: %+v", id, err)
				}
			}
		}
	}

	serviceUrl := d.Get("service_url").(string)
	subscriptionRequired := d.Get("subscription_required").(bool)

	subscriptionKeyParameterNames := expandApiManagementApiSubscriptionKeyParamNames(d.Get("subscription_key_parameter_names").([]any))

	authenticationSettings := &api.AuthenticationSettingsContract{}

	oAuth2AuthorizationSettingsRaw := d.Get("oauth2_authorization").([]any)
	authenticationSettings.OAuth2 = expandApiManagementOAuth2AuthenticationSettingsContract(oAuth2AuthorizationSettingsRaw)

	openIDAuthorizationSettingsRaw := d.Get("openid_authentication").([]any)
	authenticationSettings.Openid = expandApiManagementOpenIDAuthenticationSettingsContract(openIDAuthorizationSettingsRaw)

	contactInfoRaw := d.Get("contact").([]any)
	contactInfo := expandApiManagementApiContact(contactInfoRaw)

	licenseInfoRaw := d.Get("license").([]any)
	licenseInfo := expandApiManagementApiLicense(licenseInfoRaw)

	params := api.ApiCreateOrUpdateParameter{
		Properties: &api.ApiCreateOrUpdateProperties{
			Type:                          pointer.To(apiType),
			ApiType:                       pointer.To(soapApiType),
			Path:                          path,
			Protocols:                     protocols,
			SubscriptionKeyParameterNames: subscriptionKeyParameterNames,
			SubscriptionRequired:          &subscriptionRequired,
			AuthenticationSettings:        authenticationSettings,
			ApiRevisionDescription:        pointer.To(d.Get("revision_description").(string)),
			ApiVersionDescription:         pointer.To(d.Get("version_description").(string)),
			Contact:                       contactInfo,
			License:                       licenseInfo,
		},
	}

	if serviceUrl != "" {
		params.Properties.ServiceURL = pointer.To(serviceUrl)
	}

	if sourceApiId != "" {
		params.Properties.SourceApiId = pointer.To(sourceApiId)
	}

	if description, ok := d.GetOk("description"); ok {
		params.Properties.Description = pointer.To(description.(string))
	}

	if displayName != "" {
		params.Properties.DisplayName = pointer.To(displayName)
	}

	if version != "" {
		params.Properties.ApiVersion = pointer.To(version)
	}

	if versionSetId != "" {
		params.Properties.ApiVersionSetId = pointer.To(versionSetId)
	}

	if v, ok := d.GetOk("terms_of_service_url"); ok {
		params.Properties.TermsOfServiceURL = pointer.To(v.(string))
	}

	result, err := client.CreateOrUpdate(ctx, id, params, api.CreateOrUpdateOperationOptions{IfMatch: pointer.To("*")})
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if pollerType := custompollers.NewAPIManagementAPIPoller(client, id, result.HttpResponse); pollerType != nil {
		poller := pollers.NewPoller(pollerType, 5*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("polling creating/updating %s: %+v", id, err)
		}
	}

	return resourceApiManagementApiRead(d, meta)
}
