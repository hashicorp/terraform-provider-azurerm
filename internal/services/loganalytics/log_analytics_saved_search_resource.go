// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/savedsearches"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogAnalyticsSavedSearch() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsSavedSearchCreate,
		Read:   resourceLogAnalyticsSavedSearchRead,
		Delete: resourceLogAnalyticsSavedSearchDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := savedsearches.ParseSavedSearchID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SavedSearchV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: savedsearches.ValidateWorkspaceID,
			},

			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"category": {
				Type:         pluginsdk.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"query": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"function_alias": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"function_parameters": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					// https://learn.microsoft.com/en-us/azure/data-explorer/kusto/query/functions/user-defined-functions
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9-_]*:([a-z]+(=[^,\n]+)?|\(\*\)|(\([a-zA-Z_][a-zA-Z0-9-_]*:[a-z]+(,[a-zA-Z_][a-zA-Z0-9-_]*:([a-z]+))*\)))(,\s*[a-zA-Z_][a-zA-Z0-9-_]*:([a-z]+(=[^,\n]+)?|\(\*\)|(\([a-zA-Z_][a-zA-Z0-9-_]*:[a-z]+(,\s*[a-zA-Z_][a-zA-Z0-9-_]*:([a-z]+))*\))))*$`),
						"Log Analytics Saved Search Function Parameters must be in the following format: param-name1:type1=default_value1 OR param-name1:type1 OR param-name1:string='string goes here'",
					),
				},
			},

			"tags": tags.ForceNewSchema(),
		},
	}
}

func resourceLogAnalyticsSavedSearchCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SavedSearchesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceId, err := savedsearches.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}

	id := savedsearches.NewSavedSearchID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))
	if d.IsNewResource() {
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
			FunctionAlias: utils.String(d.Get("function_alias").(string)),
			Tags:          expandSavedSearchTag(d.Get("tags").(map[string]interface{})), // expand tags because it's defined as object set in service
		},
	}

	if v, ok := d.GetOk("function_parameters"); ok {
		attrs := v.([]interface{})
		result := make([]string, 0)
		for _, item := range attrs {
			if item != nil {
				result = append(result, item.(string))
			}
		}
		parameters.Properties.FunctionParameters = utils.String(strings.Join(result, ", "))
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsSavedSearchRead(d, meta)
}

func resourceLogAnalyticsSavedSearchRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

		functionAlias := ""
		if props.FunctionAlias != nil {
			functionAlias = *props.FunctionAlias
		}
		d.Set("function_alias", functionAlias)

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

func resourceLogAnalyticsSavedSearchDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SavedSearchesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := savedsearches.ParseSavedSearchID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandSavedSearchTag(input map[string]interface{}) *[]savedsearches.Tag {
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

func flattenSavedSearchTag(input *[]savedsearches.Tag) map[string]interface{} {
	results := make(map[string]interface{})
	if input == nil {
		return results
	}

	for _, item := range *input {
		results[item.Name] = item.Value
	}
	return results
}
