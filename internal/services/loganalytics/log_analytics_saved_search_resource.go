package loganalytics

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
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
			_, err := parse.LogAnalyticsSavedSearchID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SavedSearchV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsWorkspaceID,
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
				Type:     pluginsdk.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^[a-zA-Z0-9!-_]*:[a-zA-Z0-9!_-]+=[a-zA-Z0-9!_-]+|^[a-zA-Z0-9!-_]*:[a-zA-Z0-9!_-]+`),
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

	workspaceId, err := parse.LogAnalyticsWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewLogAnalyticsSavedSearchID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.WorkspaceName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SavedSearcheName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_log_analytics_saved_search", id.ID())
		}
	}

	parameters := operationalinsights.SavedSearch{
		SavedSearchProperties: &operationalinsights.SavedSearchProperties{
			Category:      utils.String(d.Get("category").(string)),
			DisplayName:   utils.String(d.Get("display_name").(string)),
			Query:         utils.String(d.Get("query").(string)),
			FunctionAlias: utils.String(d.Get("function_alias").(string)),
			Tags:          expandSavedSearchTag(d.Get("tags").(map[string]interface{})), // expand tags because it's defined as object set in service
		},
	}

	if v, ok := d.GetOk("function_parameters"); ok {
		attrs := v.(*pluginsdk.Set).List()
		result := make([]string, 0)
		for _, item := range attrs {
			if item != nil {
				result = append(result, item.(string))
			}
		}
		parameters.SavedSearchProperties.FunctionParameters = utils.String(strings.Join(result, ", "))
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.SavedSearcheName, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsSavedSearchRead(d, meta)
}

func resourceLogAnalyticsSavedSearchRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SavedSearchesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.LogAnalyticsSavedSearchID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.SavedSearcheName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Saved Search %q (Log Analytics Workspace %q / Resource Group %q): %s", id.WorkspaceName, id.WorkspaceName, id.ResourceGroup, err)
	}

	d.Set("name", id.SavedSearcheName)
	d.Set("log_analytics_workspace_id", parse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())

	if props := resp.SavedSearchProperties; props != nil {
		d.Set("display_name", props.DisplayName)
		d.Set("category", props.Category)
		d.Set("query", props.Query)
		d.Set("function_alias", props.FunctionAlias)
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
	id, err := parse.LogAnalyticsSavedSearchID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.SavedSearcheName); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandSavedSearchTag(input map[string]interface{}) *[]operationalinsights.Tag {
	results := make([]operationalinsights.Tag, 0)
	for key, value := range input {
		result := operationalinsights.Tag{
			Name:  utils.String(key),
			Value: utils.String(value.(string)),
		}
		results = append(results, result)
	}
	return &results
}

func flattenSavedSearchTag(input *[]operationalinsights.Tag) map[string]interface{} {
	results := make(map[string]interface{})
	if input == nil {
		return results
	}

	for _, item := range *input {
		var key string
		if item.Name != nil {
			key = *item.Name
		}
		var value string
		if item.Value != nil {
			value = *item.Value
		}
		results[key] = value
	}
	return results
}
