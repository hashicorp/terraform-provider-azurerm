// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	analyticsitems "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/analyticsitemsapis"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

var (
	userScopePath   = "myAnalyticsItems"
	sharedScopePath = "analyticsItems"
)

func resourceApplicationInsightsAnalyticsItem() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApplicationInsightsAnalyticsItemCreate,
		Read:   resourceApplicationInsightsAnalyticsItemRead,
		Update: resourceApplicationInsightsAnalyticsItemUpdate,
		Delete: resourceApplicationInsightsAnalyticsItemDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			if strings.Contains(id, userScopePath) || strings.Contains(id, sharedScopePath) {
				if _, err := analyticsitems.ParseProviderComponentID(id); err != nil {
					return err
				}
			}
			return nil
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
				ValidateFunc: components.ValidateComponentID,
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
					string(analyticsitems.ItemScopeShared),
					string(analyticsitems.ItemScopeUser),
				}, false),
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(analyticsitems.ItemTypeQuery),
					string(analyticsitems.ItemTypeFunction),
					string(analyticsitems.ItemTypeParameterFolder),
					string(analyticsitems.ItemTypeRecent),
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
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appInsightsId, err := components.ParseComponentID(d.Get("application_insights_id").(string))
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	scopeName := d.Get("scope").(string)
	typeName := d.Get("type").(string)

	itemType := analyticsitems.ItemType(typeName)
	itemScope := analyticsitems.ItemScope(scopeName)

	itemScopePath := sharedScopePath
	if itemScope == analyticsitems.ItemScopeUser {
		itemScopePath = userScopePath
	}

	id := analyticsitems.NewProviderComponentID(appInsightsId.SubscriptionId, appInsightsId.ResourceGroupName, appInsightsId.ComponentName, itemScopePath)

	// We cannot get specific analytics items without their itemID which is why we need to list all the
	// available items of a certain type and scope in order to check whether a resource already exists and needs
	// to be imported first
	// https://github.com/Azure/azure-rest-api-specs/issues/20712 itemScopePath should be set to insights.ItemScopePathAnalyticsItems in List method
	listId := analyticsitems.NewProviderComponentID(appInsightsId.SubscriptionId, appInsightsId.ResourceGroupName, appInsightsId.ComponentName, "analyticsItems")
	existing, err := client.AnalyticsItemsList(ctx, listId, analyticsitems.DefaultAnalyticsItemsListOperationOptions())
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

	properties := analyticsitems.ApplicationInsightsComponentAnalyticsItem{
		Name:    pointer.To(name),
		Type:    pointer.To(itemType),
		Scope:   pointer.To(itemScope),
		Content: pointer.To(d.Get("content").(string)),
	}
	if v := d.Get("function_alias").(string); v != "" {
		properties.Properties = &analyticsitems.ApplicationInsightsComponentAnalyticsItemProperties{
			FunctionAlias: &v,
		}
	}

	resp, err := client.AnalyticsItemsPut(ctx, id, properties, analyticsitems.DefaultAnalyticsItemsPutOperationOptions())
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if resp.Model == nil && resp.Model.Id == nil {
		return fmt.Errorf("model and model ID for %s are nil", id)
	}

	generatedId := parse.NewAnalyticsSharedItemID(id.SubscriptionId, id.ResourceGroupName, id.ComponentName, *resp.Model.Id).ID()
	if itemScope == analyticsitems.ItemScopeUser {
		generatedId = parse.NewAnalyticsUserItemID(id.SubscriptionId, id.ResourceGroupName, id.ComponentName, *resp.Model.Id).ID()
	}

	d.SetId(generatedId)

	return resourceApplicationInsightsAnalyticsItemRead(d, meta)
}

func resourceApplicationInsightsAnalyticsItemUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, itemId, err := ParseGeneratedAnalyticsItemId(d.Id())
	if err != nil {
		return err
	}

	getOptions := analyticsitems.AnalyticsItemsGetOperationOptions{
		Id: pointer.To(itemId),
	}

	existing, err := client.AnalyticsItemsGet(ctx, *id, getOptions)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	payload := existing.Model

	if d.HasChange("content") {
		payload.Content = pointer.To(d.Get("content").(string))
	}

	if d.HasChange("function_alias") {
		if payload.Properties == nil {
			payload.Properties = &analyticsitems.ApplicationInsightsComponentAnalyticsItemProperties{}
		}
		payload.Properties.FunctionAlias = pointer.To(d.Get("function_alias").(string))
	}

	putOptions := analyticsitems.AnalyticsItemsPutOperationOptions{
		OverrideItem: pointer.To(true),
	}

	if _, err = client.AnalyticsItemsPut(ctx, *id, *payload, putOptions); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceApplicationInsightsAnalyticsItemRead(d, meta)
}

func resourceApplicationInsightsAnalyticsItemRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, itemId, err := ParseGeneratedAnalyticsItemId(d.Id())
	if err != nil {
		return err
	}

	options := analyticsitems.AnalyticsItemsGetOperationOptions{
		Id: pointer.To(itemId),
	}

	resp, err := client.AnalyticsItemsGet(ctx, *id, options)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	appInsightsId := components.NewComponentID(id.SubscriptionId, id.ResourceGroupName, id.ComponentName)

	d.Set("application_insights_id", appInsightsId.ID())
	if model := resp.Model; model != nil {
		d.Set("name", model.Name)
		d.Set("version", model.Version)
		d.Set("content", model.Content)
		d.Set("scope", pointer.From(model.Scope))
		d.Set("type", pointer.From(model.Type))
		d.Set("time_created", model.TimeCreated)
		d.Set("time_modified", model.TimeModified)
		if props := model.Properties; props != nil {
			d.Set("function_alias", pointer.From(props.FunctionAlias))
		}
	}

	return nil
}

func resourceApplicationInsightsAnalyticsItemDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.AnalyticsItemsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, itemId, err := ParseGeneratedAnalyticsItemId(d.Id())
	if err != nil {
		return err
	}

	options := analyticsitems.AnalyticsItemsDeleteOperationOptions{
		Id: pointer.To(itemId),
	}

	if _, err = client.AnalyticsItemsDelete(ctx, *id, options); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func ParseGeneratedAnalyticsItemId(input string) (*analyticsitems.ProviderComponentId, string, error) {
	// The generated ID format differs depending on scope
	// <appinsightsID>/analyticsItems/<itemID>     [for shared scope items]
	// <appinsightsID>/myAnalyticsItems/<itemID>   [for user scope items]
	generatedId, err := analyticsitems.ParseProviderComponentID(input)
	if err != nil {
		return nil, "", err
	}

	scope := strings.Split(generatedId.ScopePath, "/")

	id := analyticsitems.NewProviderComponentID(generatedId.SubscriptionId, generatedId.ResourceGroupName, generatedId.ComponentName, scope[1])

	return &id, scope[2], nil
}
