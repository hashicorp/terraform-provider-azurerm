// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_api_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentapikeysapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsAPIKeyCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.APIKeysClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appInsightsId, err := componentapikeysapis.ParseComponentID(d.Get("application_insights_id").(string))
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		var existingAPIKeyList componentapikeysapis.APIKeysListOperationResponse
		var existingAPIKeyId *componentapikeysapis.ApiKeyId
		existingAPIKeyList, err = client.APIKeysList(ctx, *appInsightsId)
		if err != nil {
			if !response.WasNotFound(existingAPIKeyList.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Application Insights API key list for %s: %+v", appInsightsId, err)
			}
		}

		if existingAPIKeyList.Model != nil && len(existingAPIKeyList.Model.Value) > 0 {
			for _, existingAPIKey := range existingAPIKeyList.Model.Value {
				existingAPIKeyId, err = componentapikeysapis.ParseApiKeyIDInsensitively(*existingAPIKey.Id)
				if err != nil {
					return err
				}

				if name == *existingAPIKey.Name {
					return tf.ImportAsExistsError("azurerm_application_insights_api_key", existingAPIKeyId.ID())
				}
			}
		}
	}

	linkedReadProperties := expandApplicationInsightsAPIKeyLinkedProperties(d.Get("read_permissions").(*pluginsdk.Set), appInsightsId.ID())
	linkedWriteProperties := expandApplicationInsightsAPIKeyLinkedProperties(d.Get("write_permissions").(*pluginsdk.Set), appInsightsId.ID())
	if len(*linkedReadProperties) == 0 && len(*linkedWriteProperties) == 0 {
		return fmt.Errorf("at least one read or write permission must be defined")
	}
	apiKeyProperties := componentapikeysapis.APIKeyRequest{
		Name:                  &name,
		LinkedReadProperties:  linkedReadProperties,
		LinkedWriteProperties: linkedWriteProperties,
	}

	resp, err := client.APIKeysCreate(ctx, *appInsightsId, apiKeyProperties)
	if err != nil {
		return fmt.Errorf("creating API key %q for %s: %+v", name, appInsightsId, err)
	}

	if resp.Model == nil || resp.Model.ApiKey == nil {
		return fmt.Errorf("creating API key %q for %s: got empty API key", name, appInsightsId)
	}

	// API returns lower case on resourceGroups and apiKeys
	id, err := componentapikeysapis.ParseApiKeyIDInsensitively(*resp.Model.Id)
	if err != nil {
		return err
	}
	d.SetId(id.ID())

	// API key can only be retrieved at key creation
	d.Set("api_key", resp.Model.ApiKey)

	return resourceApplicationInsightsAPIKeyRead(d, meta)
}
