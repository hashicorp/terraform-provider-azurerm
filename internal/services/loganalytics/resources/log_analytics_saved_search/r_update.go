// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package log_analytics_saved_search

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/savedsearches"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsSavedSearchUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.SavedSearchesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := savedsearches.ParseSavedSearchID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	payload := existing.Model

	if d.HasChange("category") {
		payload.Properties.Category = d.Get("category").(string)
	}

	if d.HasChange("display_name") {
		payload.Properties.DisplayName = d.Get("display_name").(string)
	}

	if d.HasChange("query") {
		payload.Properties.Query = d.Get("query").(string)
	}

	if d.HasChange("function_alias") {
		payload.Properties.FunctionAlias = pointer.To(d.Get("function_alias").(string))
	}

	if d.HasChange("function_parameters") {
		payload.Properties.FunctionParameters = pointer.To("")
		if v := expandSavedSearchFunctionParameters(d.Get("function_parameters").([]any)); v != nil {
			payload.Properties.FunctionParameters = v
		}
	}

	if d.HasChange("tags") {
		payload.Properties.Tags = expandSavedSearchTag(d.Get("tags").(map[string]any))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceLogAnalyticsSavedSearchRead(d, meta)
}
