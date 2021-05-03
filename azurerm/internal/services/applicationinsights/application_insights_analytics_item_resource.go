package applicationinsights

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2015-05-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func resourceApplicationInsightsAnalyticsItem() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationInsightsAnalyticsItemCreate,
		Read:   resourceApplicationInsightsAnalyticsItemRead,
		Update: resourceApplicationInsightsAnalyticsItemUpdate,
		Delete: resourceApplicationInsightsAnalyticsItemDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

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

func resourceApplicationInsightsAnalyticsItemCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceApplicationInsightsAnalyticsItemCreateUpdate(d, meta, false)
}

func resourceApplicationInsightsAnalyticsItemUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceApplicationInsightsAnalyticsItemCreateUpdate(d, meta, true)
}

func resourceApplicationInsightsAnalyticsItemCreateUpdate(d *schema.ResourceData, meta interface{}, overwrite bool) error {
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appInsightsID := d.Get("application_insights_id").(string)

	resourceID, err := azure.ParseAzureResourceID(appInsightsID)
	if err != nil {
		return fmt.Errorf("Error parsing resource ID: %s", err)
	}
	resourceGroupName := resourceID.ResourceGroup
	appInsightsName := resourceID.Path["components"]

	id := d.Id()
	itemID := ""
	if id != "" {
		_, _, _, itemID, err = ResourcesArmApplicationInsightsAnalyticsItemParseID(id)
		if err != nil {
			return fmt.Errorf("Error parsing Application Insights Analytics Item ID %s: %s", id, err)
		}
	}

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

	// See comments in resourcesArmApplicationInsightsAnalyticsItemParseID method about ID format
	generatedID := appInsightsID + resourcesArmApplicationInsightsAnalyticsItemGenerateIDSuffix(itemScope, *result.ID)
	d.SetId(generatedID)

	return resourceApplicationInsightsAnalyticsItemRead(d, meta)
}

func resourceApplicationInsightsAnalyticsItemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := d.Id()
	resourceGroupName, appInsightsName, itemScopePath, itemID, err := ResourcesArmApplicationInsightsAnalyticsItemParseID(id)
	if err != nil {
		return fmt.Errorf("Error parsing Application Insights Analytics Item ID %s: %s", id, err)
	}

	result, err := client.Get(ctx, resourceGroupName, appInsightsName, itemScopePath, itemID, "")
	if err != nil {
		return fmt.Errorf("Error Getting Application Insights Analytics Item %s (Resource Group %s, App Insights Name: %s): %s", itemID, resourceGroupName, appInsightsName, err)
	}

	idSuffix := resourcesArmApplicationInsightsAnalyticsItemGenerateIDSuffix(result.Scope, itemID)
	appInsightsID := id[0 : len(id)-len(idSuffix)]
	d.Set("application_insights_id", appInsightsID)
	d.Set("name", result.Name)
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

func resourceApplicationInsightsAnalyticsItemDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := d.Id()
	resourceGroupName, appInsightsName, itemScopePath, itemID, err := ResourcesArmApplicationInsightsAnalyticsItemParseID(id)
	if err != nil {
		return fmt.Errorf("Error parsing Application Insights Analytics Item ID %s: %s", id, err)
	}

	if _, err = client.Delete(ctx, resourceGroupName, appInsightsName, itemScopePath, itemID, ""); err != nil {
		return fmt.Errorf("Error Deleting Application Insights Analytics Item '%s' (Resource Group %s, App Insights Name: %s): %s", itemID, resourceGroupName, appInsightsName, err)
	}

	return nil
}

func ResourcesArmApplicationInsightsAnalyticsItemParseID(id string) (string, string, insights.ItemScopePath, string, error) {
	resourceID, err := azure.ParseAzureResourceID(id)
	if err != nil {
		return "", "", "", "", fmt.Errorf("Error parsing resource ID: %s", err)
	}
	resourceGroupName := resourceID.ResourceGroup
	appInsightsName := resourceID.Path["components"]

	// Use the following generated ID format:
	//  <appinsightsID>/analyticsItems/<itemID>     [for shared scope items]
	//  <appinsightsID>/myanalyticsItems/<itemID>   [for user scope items]
	// Pull out the itemID and note the scope used
	itemID := resourceID.Path["analyticsItems"]
	itemScopePath := insights.AnalyticsItems
	if itemID == "" {
		// no "analyticsItems" component - try "myanalyticsItems" and set scope path
		itemID = resourceID.Path["myanalyticsItems"]
		itemScopePath = insights.MyanalyticsItems
	}

	return resourceGroupName, appInsightsName, itemScopePath, itemID, nil
}

func resourcesArmApplicationInsightsAnalyticsItemGenerateIDSuffix(itemScope insights.ItemScope, itemID string) string {
	// See comments in resourcesArmApplicationInsightsAnalyticsItemParseID method about ID format
	if itemScope == insights.ItemScopeShared {
		return "/analyticsItems/" + itemID
	} else {
		return "/myanalyticsItems/" + itemID
	}
}
