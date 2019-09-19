package azurerm

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceArmApplicationInsightsAnalyticsItem() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApplicationInsightsAnalyticsItemCreate,
		Read:   resourceArmApplicationInsightsAnalyticsItemRead,
		Update: resourceArmApplicationInsightsAnalyticsItemUpdate,
		Delete: resourceArmApplicationInsightsAnalyticsItemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"application_insights_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"content": {
				Type:     schema.TypeString,
				Required: true,
			},

			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(insights.ItemScopeShared),
					string(insights.ItemScopeUser),
				}, false),
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(insights.Query),
					string(insights.Function),
					string(insights.Folder),
					string(insights.Recent),
				}, false),
			},

			"function_alias": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"time_created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"time_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmApplicationInsightsAnalyticsItemCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmApplicationInsightsAnalyticsItemCreateUpdate(d, meta, false)
}
func resourceArmApplicationInsightsAnalyticsItemUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmApplicationInsightsAnalyticsItemCreateUpdate(d, meta, true)
}
func resourceArmApplicationInsightsAnalyticsItemCreateUpdate(d *schema.ResourceData, meta interface{}, overwrite bool) error {
	client := meta.(*ArmClient).appInsights.AnalyticsItemsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	appInsightsID := d.Get("application_insights_id").(string)

	id, err := azure.ParseAzureResourceID(appInsightsID)
	if err != nil {
		return fmt.Errorf("Error parsing resource ID: %s", err)
	}

	appInsightsName := id.Path["components"]

	itemID := d.Id()
	name := d.Get("name").(string)
	content := d.Get("content").(string)
	scopeName := d.Get("scope").(string)
	typeName := d.Get("type").(string)
	functionAlias := d.Get("function_alias").(string)

	itemType := insights.ItemType(typeName)
	itemScope := insights.ItemScope(scopeName)
	properties := insights.ApplicationInsightsComponentAnalyticsItem{
		ID:      &itemID,
		Name:    &name,
		Type:    itemType,
		Scope:   itemScope,
		Content: &content,
	}
	if functionAlias != "" {
		properties.Properties = &insights.ApplicationInsightsComponentAnalyticsItemProperties{
			FunctionAlias: &functionAlias,
		}
	}

	var itemScopePath insights.ItemScopePath
	if itemScope == insights.ItemScopeUser {
		itemScopePath = insights.MyanalyticsItems
	} else {
		itemScopePath = insights.AnalyticsItems
	}
	result, err := client.Put(ctx, resourceGroupName, appInsightsName, itemScopePath, properties, &overwrite)
	if err != nil {
		return fmt.Errorf("Error Putting Application Insights Analytics Item %s (Resource Group %s, App Insights Name: %s): %s", name, resourceGroupName, appInsightsName, err)
	}

	d.SetId(*result.ID)

	return resourceArmApplicationInsightsAnalyticsItemRead(d, meta)
}

func resourceArmApplicationInsightsAnalyticsItemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appInsights.AnalyticsItemsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	appInsightsID := d.Get("application_insights_id").(string)
	scopeName := d.Get("scope").(string)
	itemID := d.Id()

	id, err := azure.ParseAzureResourceID(appInsightsID)
	if err != nil {
		return fmt.Errorf("Error parsing resource ID: %s", err)
	}

	appInsightsName := id.Path["components"]
	name := d.Get("name").(string)

	var scopePath insights.ItemScopePath
	if scopeName == "user" {
		scopePath = insights.MyanalyticsItems
	} else {
		scopePath = insights.AnalyticsItems
	}
	result, err := client.Get(ctx, resourceGroupName, appInsightsName, scopePath, itemID, name)
	if err != nil {
		return fmt.Errorf("Error Getting Application Insights Analytics Item %s (Resource Group %s, App Insights Name: %s): %s", name, resourceGroupName, appInsightsName, err)
	}

	d.Set("version", result.Version)
	d.Set("content", result.Content)
	d.Set("scope", string(result.Scope))
	d.Set("type", string(result.Type))
	d.Set("time_created", result.TimeCreated)
	d.Set("time_modified", result.TimeModified)

	if result.Properties != nil {
		d.Set("function_alias", result.Properties.FunctionAlias)
	}

	return nil
}

func resourceArmApplicationInsightsAnalyticsItemDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ArmClient).appInsights.AnalyticsItemsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	appInsightsID := d.Get("application_insights_id").(string)
	scopeName := d.Get("scope").(string)
	itemID := d.Id()

	id, err := azure.ParseAzureResourceID(appInsightsID)
	if err != nil {
		return fmt.Errorf("Error parsing resource ID: %s", err)
	}

	appInsightsName := id.Path["components"]
	name := d.Get("name").(string)

	var scopePath insights.ItemScopePath
	if scopeName == string(insights.ItemScopeUser) {
		scopePath = insights.MyanalyticsItems
	} else {
		scopePath = insights.AnalyticsItems
	}
	_, err = client.Delete(ctx, resourceGroupName, appInsightsName, scopePath, itemID, name)
	if err != nil {
		return fmt.Errorf("Error Getting Application Insights Analytics Item %s (Resource Group %s, App Insights Name: %s): %s", name, resourceGroupName, appInsightsName, err)
	}

	return nil
}
