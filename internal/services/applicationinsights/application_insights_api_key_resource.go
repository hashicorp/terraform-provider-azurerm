// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	apikeys "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentapikeysapis"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsAPIKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApplicationInsightsAPIKeyCreate,
		Read:   resourceApplicationInsightsAPIKeyRead,
		Delete: resourceApplicationInsightsAPIKeyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apikeys.ParseApiKeyID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiKeyUpgradeV0ToV1{},
			1: migration.ApiKeyUpgradeV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"application_insights_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: components.ValidateComponentID,
			},

			"read_permissions": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				ForceNew: true,
				Set:      pluginsdk.HashString,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"agentconfig", "aggregate", "api", "draft", "extendqueries", "search"}, false),
				},
			},

			"write_permissions": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				ForceNew: true,
				Set:      pluginsdk.HashString,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"annotations"}, false),
				},
			},

			"api_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceApplicationInsightsAPIKeyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.APIKeysClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appInsightsId, err := apikeys.ParseComponentID(d.Get("application_insights_id").(string))
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	var existingAPIKeyList apikeys.APIKeysListOperationResponse
	var existingAPIKeyId *apikeys.ApiKeyId
	existingAPIKeyList, err = client.APIKeysList(ctx, *appInsightsId)
	if err != nil {
		if !response.WasNotFound(existingAPIKeyList.HttpResponse) {
			return fmt.Errorf("checking for presence of existing Application Insights API key list for %s: %+v", appInsightsId, err)
		}
	}

	for existingAPIKeyList.Model != nil && len(existingAPIKeyList.Model.Value) > 0 {
		for _, existingAPIKey := range existingAPIKeyList.Model.Value {
			existingAPIKeyId, err = apikeys.ParseApiKeyIDInsensitively(*existingAPIKey.Id)
			if err != nil {
				return err
			}

			if name == *existingAPIKey.Name {
				return tf.ImportAsExistsError("azurerm_application_insights_api_key", existingAPIKeyId.ID())
			}
		}
	}

	linkedReadProperties := expandApplicationInsightsAPIKeyLinkedProperties(d.Get("read_permissions").(*pluginsdk.Set), appInsightsId.ID())
	linkedWriteProperties := expandApplicationInsightsAPIKeyLinkedProperties(d.Get("write_permissions").(*pluginsdk.Set), appInsightsId.ID())
	if len(*linkedReadProperties) == 0 && len(*linkedWriteProperties) == 0 {
		return fmt.Errorf("at least one read or write permission must be defined")
	}
	apiKeyProperties := apikeys.APIKeyRequest{
		Name:                  &name,
		LinkedReadProperties:  linkedReadProperties,
		LinkedWriteProperties: linkedWriteProperties,
	}

	resp, err := client.APIKeysCreate(ctx, *appInsightsId, apiKeyProperties)
	if err != nil {
		return fmt.Errorf("creating API key %q for %s: %+v", name, appInsightsId, err)
	}

	if resp.Model == nil || resp.Model.ApiKey == nil {
		return fmt.Errorf("creating API key %q for %s: got empty API key", name, appInsightsId)
	}

	// API returns lower case on resourceGroups and apiKeys
	id, err := apikeys.ParseApiKeyIDInsensitively(*resp.Model.Id)
	if err != nil {
		return err
	}
	d.SetId(id.ID())

	// API key can only be retrieved at key creation
	d.Set("api_key", resp.Model.ApiKey)

	return resourceApplicationInsightsAPIKeyRead(d, meta)
}

func resourceApplicationInsightsAPIKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.APIKeysClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apikeys.ParseApiKeyID(d.Id())
	if err != nil {
		return err
	}

	appInsightsId := components.NewComponentID(subscriptionId, id.ResourceGroupName, id.ComponentName)

	result, err := client.APIKeysGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			log.Printf("[DEBUG] %s not found, removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("application_insights_id", appInsightsId.ID())

	if model := result.Model; model != nil {
		d.Set("name", model.Name)
		if props := model.LinkedReadProperties; props != nil {
			readProps := flattenApplicationInsightsAPIKeyLinkedProperties(props)
			if err := d.Set("read_permissions", readProps); err != nil {
				return fmt.Errorf("flattening `read_permissions `: %s", err)
			}
		}
		if props := model.LinkedWriteProperties; props != nil {
			writeProps := flattenApplicationInsightsAPIKeyLinkedProperties(props)
			if err := d.Set("write_permissions", writeProps); err != nil {
				return fmt.Errorf("flattening `write_permissions `: %s", err)
			}
		}
	}

	return nil
}

func resourceApplicationInsightsAPIKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.APIKeysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apikeys.ParseApiKeyID(d.Id())
	if err != nil {
		return err
	}

	result, err := client.APIKeysDelete(ctx, *id)
	if err != nil {
		if response.WasNotFound(result.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
