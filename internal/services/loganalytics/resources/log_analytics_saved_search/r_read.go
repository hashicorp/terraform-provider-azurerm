// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_saved_search

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/savedsearches"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsSavedSearchRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.SavedSearchesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := savedsearches.ParseSavedSearchID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	d.Set("name", id.SavedSearchId)
	d.Set("log_analytics_workspace_id", savedsearches.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())

	if model := resp.Model; model != nil {
		props := model.Properties

		d.Set("display_name", props.DisplayName)
		d.Set("category", props.Category)
		d.Set("query", props.Query)

		d.Set("function_alias", pointer.From(props.FunctionAlias))

		functionParams := make([]string, 0)
		if props.FunctionParameters != nil {
			functionParams = strings.Split(*props.FunctionParameters, ", ")
		}
		d.Set("function_parameters", functionParams)

		// flatten tags because it's defined as object set in service
		if err := d.Set("tags", flattenSavedSearchTag(props.Tags)); err != nil {
			return fmt.Errorf("setting `tag`: %+v", err)
		}
	}

	return nil
}

func flattenSavedSearchTag(input *[]savedsearches.Tag) map[string]any {
	results := make(map[string]any)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results[item.Name] = item.Value
	}
	return results
}
