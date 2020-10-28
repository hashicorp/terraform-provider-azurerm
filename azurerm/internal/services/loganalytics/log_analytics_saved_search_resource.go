package loganalytics

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2020-03-01-preview/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogAnalyticsSavedSearch() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogAnalyticsSavedSearchCreate,
		Read:   resourceArmLogAnalyticsSavedSearchRead,
		Delete: resourceArmLogAnalyticsSavedSearchDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"log_analytics_workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.LogAnalyticsWorkspaceID,
				// https://github.com/Azure/azure-rest-api-specs/issues/9330
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// https://github.com/Azure/azure-rest-api-specs/issues/9330
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"category": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"query": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"function_alias": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"function_parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^[a-zA-Z0-9!-_]*:[a-zA-Z0-9!_-]+=[a-zA-Z0-9!_-]+`),
						"Log Analytics Saved Search Function Parameters must be in the following format: param-name1:type1=default_value1",
					),
				},
			},

			"tags": tags.ForceNewSchema(),
		},
	}
}

func resourceArmLogAnalyticsSavedSearchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SavedSearchesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Saved Search creation.")

	name := d.Get("name").(string)
	workspaceID := d.Get("log_analytics_workspace_id").(string)
	id, err := parse.LogAnalyticsWorkspaceID(workspaceID)
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Log Analytics Saved Search %q (WorkSpace %q / Resource Group %q): %s", name, id.Name, id.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_saved_search", *existing.ID)
		}
	}

	parameters := operationalinsights.SavedSearch{
		SavedSearchProperties: &operationalinsights.SavedSearchProperties{
			Category:      utils.String(d.Get("category").(string)),
			DisplayName:   utils.String(d.Get("display_name").(string)),
			Query:         utils.String(d.Get("query").(string)),
			FunctionAlias: utils.String(d.Get("function_alias").(string)),
			Tags:          expandArmSavedSearchTag(d.Get("tags").(map[string]interface{})), // expand tags because it's defined as object set in service
		},
	}

	if v, ok := d.GetOk("function_parameters"); ok {
		attrs := v.(*schema.Set).List()
		result := make([]string, 0)
		for _, item := range attrs {
			if item != nil {
				result = append(result, item.(string))
			}
		}
		parameters.SavedSearchProperties.FunctionParameters = utils.String(strings.Join(result, ", "))
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, name, parameters); err != nil {
		return fmt.Errorf("creating Saved Search %q (Log Analytics Workspace %q / Resource Group %q): %+v", name, id.Name, id.ResourceGroup, err)
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Saved Search %q (Log Analytics Workspace %q / Resource Group %q): %+v", name, id.Name, id.ResourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read Log Analytics Saved Search %q (WorkSpace %q / Resource Group %q): %s", name, id.Name, id.ResourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceArmLogAnalyticsSavedSearchRead(d, meta)
}

func resourceArmLogAnalyticsSavedSearchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SavedSearchesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.LogAnalyticsSavedSearchID(d.Id())
	if err != nil {
		return err
	}
	workspaceId := parse.NewLogAnalyticsWorkspaceID(id.WorkspaceName, id.ResourceGroup).ID(meta.(*clients.Client).Account.SubscriptionId)

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Saved Search %q (Log Analytics Workspace %q / Resource Group %q): %s", id.Name, id.WorkspaceName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("log_analytics_workspace_id", workspaceId)

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
		if err := d.Set("tags", flattenArmSavedSearchTag(props.Tags)); err != nil {
			return fmt.Errorf("setting `tag`: %+v", err)
		}
	}

	return nil
}

func resourceArmLogAnalyticsSavedSearchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SavedSearchesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.LogAnalyticsSavedSearchID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
		return fmt.Errorf("deleting Saved Search %q (Log Analytics Workspace %q / Resource Group %q): %s", id.Name, id.WorkspaceName, id.ResourceGroup, err)
	}

	return nil
}

func expandArmSavedSearchTag(input map[string]interface{}) *[]operationalinsights.Tag {
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

func flattenArmSavedSearchTag(input *[]operationalinsights.Tag) map[string]interface{} {
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
