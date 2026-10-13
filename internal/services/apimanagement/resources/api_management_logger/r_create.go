// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_logger

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementLoggerCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := logger.NewLoggerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	eventHubRaw := d.Get("eventhub").([]any)
	appInsightsRaw := d.Get("application_insights").([]any)

	if len(eventHubRaw) == 0 && len(appInsightsRaw) == 0 {
		return errors.New("either `eventhub` or `application_insights` is required")
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_logger", id.ID())
			}
		}
	}

	parameters := logger.LoggerContract{
		Properties: &logger.LoggerContractProperties{
			IsBuffered:  pointer.To(d.Get("buffered").(bool)),
			Description: pointer.To(d.Get("description").(string)),
		},
	}

	if len(eventHubRaw) > 0 {
		parameters.Properties.LoggerType = logger.LoggerTypeAzureEventHub
		parameters.Properties.Credentials = expandApiManagementLoggerEventHub(eventHubRaw)
	} else if len(appInsightsRaw) > 0 {
		parameters.Properties.LoggerType = logger.LoggerTypeApplicationInsights
		parameters.Properties.Credentials = expandApiManagementLoggerApplicationInsights(appInsightsRaw)
	}

	if resourceId := d.Get("resource_id").(string); resourceId != "" {
		parameters.Properties.ResourceId = pointer.To(resourceId)
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, logger.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementLoggerRead(d, meta)
}

func expandApiManagementLoggerEventHub(input []any) *map[string]string {
	credentials := make(map[string]string)
	eventHub := input[0].(map[string]any)

	connectionString := eventHub["connection_string"].(string)
	endpointAddress := eventHub["endpoint_uri"].(string)
	clientId := eventHub["user_assigned_identity_client_id"].(string)

	credentials["name"] = eventHub["name"].(string)
	if len(connectionString) > 0 {
		credentials["connectionString"] = connectionString
	} else if len(endpointAddress) > 0 {
		credentials["endpointAddress"] = endpointAddress
		// This field is required by the API and only accepts either a valid UUID or `SystemAssigned` as a value, so we default this to `SystemAssigned` in the create if the field is omitted
		credentials["identityClientId"] = "SystemAssigned"
		if clientId != "" {
			credentials["identityClientId"] = clientId
		}
	}

	return &credentials
}

func expandApiManagementLoggerApplicationInsights(input []any) *map[string]string {
	credentials := make(map[string]string)
	ai := input[0].(map[string]any)
	if ai["instrumentation_key"].(string) != "" {
		credentials["instrumentationKey"] = ai["instrumentation_key"].(string)
	}
	if ai["connection_string"].(string) != "" {
		credentials["connectionString"] = ai["connection_string"].(string)
	}
	if clientId := ai["identity_client_id"].(string); clientId != "" {
		credentials["identityClientId"] = clientId
	}
	return &credentials
}
