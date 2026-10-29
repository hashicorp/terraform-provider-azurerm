// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_saved_search

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/savedsearches"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsSavedSearchCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.SavedSearchesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := savedsearches.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}

	id := savedsearches.NewSavedSearchID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_log_analytics_saved_search", id.ID())
		}
	}

	parameters := savedsearches.SavedSearch{
		Properties: savedsearches.SavedSearchProperties{
			Category:      d.Get("category").(string),
			DisplayName:   d.Get("display_name").(string),
			Query:         d.Get("query").(string),
			FunctionAlias: pointer.To(d.Get("function_alias").(string)),
			Tags:          expandSavedSearchTag(d.Get("tags").(map[string]any)), // expand tags because it's defined as object set in service
		},
	}

	if v, ok := d.GetOk("function_parameters"); ok {
		attrs := v.([]any)
		result := make([]string, 0)
		for _, item := range attrs {
			if item != nil {
				result = append(result, item.(string))
			}
		}
		parameters.Properties.FunctionParameters = pointer.To(strings.Join(result, ", "))
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsSavedSearchRead(d, meta)
}

func expandSavedSearchTag(input map[string]any) *[]savedsearches.Tag {
	results := make([]savedsearches.Tag, 0)
	for key, value := range input {
		result := savedsearches.Tag{
			Name:  key,
			Value: value.(string),
		}
		results = append(results, result)
	}
	return &results
}
