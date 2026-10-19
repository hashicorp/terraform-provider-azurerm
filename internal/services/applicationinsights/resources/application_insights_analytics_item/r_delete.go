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

func resourceApplicationInsightsAnalyticsItemDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, itemId, err := ParseGeneratedAnalyticsItemId(d.Id())
	if err != nil {
		return err
	}

	options := analyticsitemsapis.AnalyticsItemsDeleteOperationOptions{
		Id: pointer.To(itemId),
	}

	if _, err = client.AnalyticsItemsDelete(ctx, *id, options); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
