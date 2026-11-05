// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_analytics_item

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/analyticsitemsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsAnalyticsItemCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appInsightsId, err := componentsapis.ParseComponentID(d.Get("application_insights_id").(string))
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	scopeName := d.Get("scope").(string)
	typeName := d.Get("type").(string)

	itemType := analyticsitemsapis.ItemType(typeName)
	itemScope := analyticsitemsapis.ItemScope(scopeName)

	itemScopePath := sharedScopePath
	if itemScope == analyticsitemsapis.ItemScopeUser {
		itemScopePath = userScopePath
	}

	id := analyticsitemsapis.NewProviderComponentID(appInsightsId.SubscriptionId, appInsightsId.ResourceGroupName, appInsightsId.ComponentName, itemScopePath)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		// We cannot get specific analytics items without their itemID which is why we need to list all the
		// available items of a certain type and scope in order to check whether a resource already exists and needs
		// to be imported first
		// https://github.com/Azure/azure-rest-api-specs/issues/20712 itemScopePath should be set to insights.ItemScopePathAnalyticsItems in List method
		listId := analyticsitemsapis.NewProviderComponentID(appInsightsId.SubscriptionId, appInsightsId.ResourceGroupName, appInsightsId.ComponentName, "analyticsItems")
		existing, err := client.AnalyticsItemsList(ctx, listId, analyticsitemsapis.DefaultAnalyticsItemsListOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if model := existing.Model; model != nil {
			for _, value := range *model {
				if v := value.Name; v != nil && *v == name {
					return tf.ImportAsExistsError("azurerm_application_insights_analytics_item", *value.Id)
				}
			}
		}
	}

	properties := analyticsitemsapis.ApplicationInsightsComponentAnalyticsItem{
		Name:    pointer.To(name),
		Type:    pointer.To(itemType),
		Scope:   pointer.To(itemScope),
		Content: pointer.To(d.Get("content").(string)),
	}
	if v := d.Get("function_alias").(string); v != "" {
		properties.Properties = &analyticsitemsapis.ApplicationInsightsComponentAnalyticsItemProperties{
			FunctionAlias: &v,
		}
	}

	resp, err := client.AnalyticsItemsPut(ctx, id, properties, analyticsitemsapis.DefaultAnalyticsItemsPutOperationOptions())
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if resp.Model == nil && resp.Model.Id == nil {
		return fmt.Errorf("model and model ID for %s are nil", id)
	}

	generatedId := parse.NewAnalyticsSharedItemID(id.SubscriptionId, id.ResourceGroupName, id.ComponentName, *resp.Model.Id).ID()
	if itemScope == analyticsitemsapis.ItemScopeUser {
		generatedId = parse.NewAnalyticsUserItemID(id.SubscriptionId, id.ResourceGroupName, id.ComponentName, *resp.Model.Id).ID()
	}

	d.SetId(generatedId)

	return resourceApplicationInsightsAnalyticsItemRead(d, meta)
}
