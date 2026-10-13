// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_logger

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementLoggerUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := logger.NewLoggerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	eventHubRaw, hasEventHub := d.GetOk("eventhub")
	appInsightsRaw, hasAppInsights := d.GetOk("application_insights")

	parameters := logger.LoggerUpdateContract{
		Properties: &logger.LoggerUpdateParameters{
			IsBuffered:  pointer.To(d.Get("buffered").(bool)),
			Description: pointer.To(d.Get("description").(string)),
		},
	}

	if hasEventHub {
		parameters.Properties.LoggerType = pointer.To(logger.LoggerTypeAzureEventHub)
		parameters.Properties.Credentials = expandApiManagementLoggerEventHub(eventHubRaw.([]any))
	} else if hasAppInsights {
		parameters.Properties.LoggerType = pointer.To(logger.LoggerTypeApplicationInsights)
		parameters.Properties.Credentials = expandApiManagementLoggerApplicationInsights(appInsightsRaw.([]any))
	}

	if _, err := client.Update(ctx, id, parameters, logger.UpdateOperationOptions{}); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceApiManagementLoggerRead(d, meta)
}
