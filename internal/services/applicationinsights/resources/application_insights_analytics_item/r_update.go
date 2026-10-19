// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_analytics_item

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/analyticsitemsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsAnalyticsItemUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, itemId, err := ParseGeneratedAnalyticsItemId(d.Id())
	if err != nil {
		return err
	}

	getOptions := analyticsitemsapis.AnalyticsItemsGetOperationOptions{
		Id: pointer.To(itemId),
	}

	existing, err := client.AnalyticsItemsGet(ctx, *id, getOptions)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	payload := existing.Model

	if d.HasChange("content") {
		payload.Content = pointer.To(d.Get("content").(string))
	}

	if d.HasChange("function_alias") {
		if payload.Properties == nil {
			payload.Properties = &analyticsitemsapis.ApplicationInsightsComponentAnalyticsItemProperties{}
		}
		payload.Properties.FunctionAlias = pointer.To(d.Get("function_alias").(string))
	}

	putOptions := analyticsitemsapis.AnalyticsItemsPutOperationOptions{
		OverrideItem: pointer.To(true),
	}

	if _, err = client.AnalyticsItemsPut(ctx, *id, *payload, putOptions); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceApplicationInsightsAnalyticsItemRead(d, meta)
}
