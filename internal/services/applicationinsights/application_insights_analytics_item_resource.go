// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2020-02-02/insights" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApplicationInsightsAnalyticsItem() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApplicationInsightsAnalyticsItemCreate,
		Read:   resourceApplicationInsightsAnalyticsItemRead,
		Update: resourceApplicationInsightsAnalyticsItemUpdate,
		Delete: resourceApplicationInsightsAnalyticsItemDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			if strings.Contains(id, "myAnalyticsItems") {
				_, err := parse.AnalyticsUserItemID(id)
				return err
			} else {
				strings.Contains(id, string(insights.ItemScopePathAnalyticsItems))
				_, err := parse.AnalyticsSharedItemID(id)
				return err
			}
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AnalyticsItemUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"application_insights_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ComponentID,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"content": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"scope": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(insights.ItemScopeShared),
					string(insights.ItemScopeUser),
				}, false),
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(insights.ItemTypeQuery),
					string(insights.ItemTypeFunction),
					string(insights.ItemTypeFolder),
					string(insights.ItemTypeRecent),
				}, false),
			},

			"function_alias": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"time_created": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"time_modified": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApplicationInsightsAnalyticsItemCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceApplicationInsightsAnalyticsItemCreateUpdate(d, meta, false)
}

func resourceApplicationInsightsAnalyticsItemUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	return resourceApplicationInsightsAnalyticsItemCreateUpdate(d, meta, true)
}

func resourceApplicationInsightsAnalyticsItemCreateUpdate(d *pluginsdk.ResourceData, meta interface{}, overwrite bool) error {
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appInsightsId, err := parse.ComponentID(d.Get("application_insights_id").(string))
	if err != nil {
		return err
	}

	var itemID string
	var id string
	if id, _, _, _, itemID, err = ResourcesArmApplicationInsightsAnalyticsItemParseID(d.Id()); d.Id() != "" {
		if err != nil {
			return fmt.Errorf("parsing Application Insights Analytics Item ID %s: %+v", d.Id(), err)
		}
	}

	name := d.Get("name").(string)
	content := d.Get("content").(string)
	scopeName := d.Get("scope").(string)
	typeName := d.Get("type").(string)
	functionAlias := d.Get("function_alias").(string)

	itemType := insights.ItemType(typeName)
	itemScope := insights.ItemScope(scopeName)

	var itemScopePath insights.ItemScopePath
	if itemScope == insights.ItemScopeUser {
		itemScopePath = insights.ItemScopePathMyanalyticsItems
	} else {
		itemScopePath = insights.ItemScopePathAnalyticsItems
	}

	includeContent := false

	if d.IsNewResource() {
		// We cannot get specific analytics items without their itemID which is why we need to list all the
		// available items of a certain type and scope in order to check whether a resource already exists and needs
		// to be imported first
		// https://github.com/Azure/azure-rest-api-specs/issues/20712 itemScopePath should be set to insights.ItemScopePathAnalyticsItems in List method
		existing, err := client.List(ctx, appInsightsId.ResourceGroup, appInsightsId.Name, insights.ItemScopePathAnalyticsItems, itemScope, insights.ItemTypeParameter(typeName), &includeContent)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Application Insights Analytics Items %+v", err)
			}
		}
		if existing.Value != nil {
			values := *existing.Value
			for _, v := range values {
				if *v.Name == name {
					return tf.ImportAsExistsError("azurerm_application_insights_analytics_item", *v.ID)
				}
			}
		}
	}

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

	result, err := client.Put(ctx, appInsightsId.ResourceGroup, appInsightsId.Name, itemScopePath, properties, &overwrite)
	if err != nil {
		return fmt.Errorf("putting Application Insights Analytics Item %s: %+v", id, err)
	}

	// See comments in ResourcesArmApplicationInsightsAnalyticsItemParseID method about ID format
	generatedID := resourcesArmApplicationInsightsAnalyticsItemGenerateID(itemScope, *result.ID, appInsightsId)

	d.SetId(generatedID)

	return resourceApplicationInsightsAnalyticsItemRead(d, meta)
}

func resourceApplicationInsightsAnalyticsItemRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	subscriptionId := meta.(*clients.Client).AppInsights.AnalyticsItemsClient.SubscriptionID
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, resourceGroupName, appInsightsName, itemScopePath, itemID, err := ResourcesArmApplicationInsightsAnalyticsItemParseID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Application Insights Analytics Item ID %s: %s", d.Id(), err)
	}

	result, err := client.Get(ctx, resourceGroupName, appInsightsName, itemScopePath, itemID, "")
	if err != nil {
		if utils.ResponseWasNotFound(result.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("getting Application Insights Analytics Item %s: %+v", id, err)
	}

	appInsightsId := parse.NewComponentID(subscriptionId, resourceGroupName, appInsightsName)

	d.Set("application_insights_id", appInsightsId.ID())
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

func resourceApplicationInsightsAnalyticsItemDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, resourceGroupName, appInsightsName, itemScopePath, itemID, err := ResourcesArmApplicationInsightsAnalyticsItemParseID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Application Insights Analytics Item ID %s: %+v", d.Id(), err)
	}

	if _, err = client.Delete(ctx, resourceGroupName, appInsightsName, itemScopePath, itemID, ""); err != nil {
		return fmt.Errorf("deleting Application Insights Analytics Item %s: %+v", id, err)
	}

	return nil
}

func ResourcesArmApplicationInsightsAnalyticsItemParseID(id string) (string, string, string, insights.ItemScopePath, string, error) {
	// The generated ID format differs depending on scope
	// <appinsightsID>/analyticsItems/<itemID>     [for shared scope items]
	// <appinsightsID>/myAnalyticsItems/<itemID>   [for user scope items]
	switch {
	case strings.Contains(id, "myAnalyticsItems"):
		id, err := parse.AnalyticsUserItemID(id)
		if err != nil {
			return "", "", "", "", "", err
		}
		return id.String(), id.ResourceGroup, id.ComponentName, insights.ItemScopePathMyanalyticsItems, id.MyAnalyticsItemName, nil
	case strings.Contains(id, string(insights.ItemScopePathAnalyticsItems)):
		id, err := parse.AnalyticsSharedItemID(id)
		if err != nil {
			return "", "", "", "", "", err
		}
		return id.String(), id.ResourceGroup, id.ComponentName, insights.ItemScopePathAnalyticsItems, id.AnalyticsItemName, nil
	default:
		return "", "", "", "", "", fmt.Errorf("parsing Application Insights Analytics Item ID %s", id)
	}
}

func resourcesArmApplicationInsightsAnalyticsItemGenerateID(itemScope insights.ItemScope, itemID string, appInsightsId *parse.ComponentId) string {
	if itemScope == insights.ItemScopeShared {
		id := parse.NewAnalyticsSharedItemID(appInsightsId.SubscriptionId, appInsightsId.ResourceGroup, appInsightsId.Name, itemID)
		return id.ID()
	} else {
		id := parse.NewAnalyticsUserItemID(appInsightsId.SubscriptionId, appInsightsId.ResourceGroup, appInsightsId.Name, itemID)
		return id.ID()
	}
}
