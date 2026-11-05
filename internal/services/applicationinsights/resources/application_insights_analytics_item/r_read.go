// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_analytics_item

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/analyticsitemsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsAnalyticsItemRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, itemId, err := ParseGeneratedAnalyticsItemId(d.Id())
	if err != nil {
		return err
	}

	options := analyticsitemsapis.AnalyticsItemsGetOperationOptions{
		Id: pointer.To(itemId),
	}

	resp, err := client.AnalyticsItemsGet(ctx, *id, options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	appInsightsId := componentsapis.NewComponentID(id.SubscriptionId, id.ResourceGroupName, id.ComponentName)

	d.Set("application_insights_id", appInsightsId.ID())
	if model := resp.Model; model != nil {
		d.Set("name", model.Name)
		d.Set("version", model.Version)
		d.Set("content", model.Content)
		d.Set("scope", pointer.From(model.Scope))
		d.Set("type", pointer.From(model.Type))
		d.Set("time_created", model.TimeCreated)
		d.Set("time_modified", model.TimeModified)
		if props := model.Properties; props != nil {
			d.Set("function_alias", pointer.From(props.FunctionAlias))
		}
	}

	return nil
}

var (
	userScopePath   = "myAnalyticsItems"
	sharedScopePath = "analyticsItems"
)

func ParseGeneratedAnalyticsItemId(input string) (*analyticsitemsapis.ProviderComponentId, string, error) {
	// The generated ID format differs depending on scope
	// <appinsightsID>/analyticsItems/<itemID>     [for shared scope items]
	// <appinsightsID>/myAnalyticsItems/<itemID>   [for user scope items]
	generatedId, err := analyticsitemsapis.ParseProviderComponentID(input)
	if err != nil {
		return nil, "", err
	}

	scope := strings.Split(generatedId.ScopePath, "/")

	id := analyticsitemsapis.NewProviderComponentID(generatedId.SubscriptionId, generatedId.ResourceGroupName, generatedId.ComponentName, scope[1])

	return &id, scope[2], nil
}
