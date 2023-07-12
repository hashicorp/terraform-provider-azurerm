// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"fmt"
	"log"
	"regexp"
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

func resourceApplicationInsightsAPIKey() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApplicationInsightsAPIKeyCreate,
		Read:   resourceApplicationInsightsAPIKeyRead,
		Delete: resourceApplicationInsightsAPIKeyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ApiKeyID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiKeyUpgradeV0ToV1{},
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
				ValidateFunc: validation.NoZeroValues,
			},

			"application_insights_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ComponentID,
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

	log.Printf("[INFO] preparing arguments for AzureRM Application Insights API key creation.")

	appInsightsId, err := parse.ComponentID(d.Get("application_insights_id").(string))
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	var existingAPIKeyList insights.ApplicationInsightsComponentAPIKeyListResult
	var existingAPIKeyId *parse.ApiKeyId
	var keyId string
	existingAPIKeyList, err = client.List(ctx, appInsightsId.ResourceGroup, appInsightsId.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existingAPIKeyList.Response) {
			return fmt.Errorf("checking for presence of existing Application Insights API key list %q (%s): %+v", name, appInsightsId, err)
		}
	}

	if existingAPIKeyList.Value != nil {
		for _, existingAPIKey := range *existingAPIKeyList.Value {
			existingAPIKeyId, err = parse.ApiKeyID(camelCaseApiKeys(*existingAPIKey.ID))
			if err != nil {
				return err
			}

			existingAppInsightsName := existingAPIKeyId.ComponentName
			if appInsightsId.Name == existingAppInsightsName {
				keyId = existingAPIKeyId.Name
				break
			}
		}
	}

	var existing insights.ApplicationInsightsComponentAPIKey
	existing, err = client.Get(ctx, appInsightsId.ResourceGroup, appInsightsId.Name, keyId)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Application Insights API key %q (%s): %s", name, appInsightsId, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) && existingAPIKeyId != nil {
		return tf.ImportAsExistsError("azurerm_application_insights_api_key", existingAPIKeyId.ID())
	}

	linkedReadProperties := expandApplicationInsightsAPIKeyLinkedProperties(d.Get("read_permissions").(*pluginsdk.Set), appInsightsId.ID())
	linkedWriteProperties := expandApplicationInsightsAPIKeyLinkedProperties(d.Get("write_permissions").(*pluginsdk.Set), appInsightsId.ID())
	if len(*linkedReadProperties) == 0 && len(*linkedWriteProperties) == 0 {
		return fmt.Errorf("at least one read or write permission must be defined")
	}
	apiKeyProperties := insights.APIKeyRequest{
		Name:                  &name,
		LinkedReadProperties:  linkedReadProperties,
		LinkedWriteProperties: linkedWriteProperties,
	}

	result, err := client.Create(ctx, appInsightsId.ResourceGroup, appInsightsId.Name, apiKeyProperties)
	if err != nil {
		return fmt.Errorf("creating Application Insights API key %q (%s): %+v", name, appInsightsId, err)
	}

	if result.APIKey == nil {
		return fmt.Errorf("creating Application Insights API key %q (%s): got empty API key", name, appInsightsId)
	}

	d.SetId(camelCaseApiKeys(*result.ID))

	// API key can only retrieved at key creation
	d.Set("api_key", result.APIKey)

	return resourceApplicationInsightsAPIKeyRead(d, meta)
}

func resourceApplicationInsightsAPIKeyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.APIKeysClient
	subscriptionId := meta.(*clients.Client).AppInsights.APIKeysClient.SubscriptionID
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiKeyID(d.Id())
	if err != nil {
		return err
	}

	appInsightsId := parse.NewComponentID(subscriptionId, id.ResourceGroup, id.ComponentName)

	log.Printf("[DEBUG] Reading AzureRM Application Insights API key '%s'", id)

	result, err := client.Get(ctx, id.ResourceGroup, id.ComponentName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(result.Response) {
			log.Printf("[WARN] AzureRM Application Insights API key '%s' not found, removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.Set("application_insights_id", appInsightsId.ID())

	d.Set("name", result.Name)
	readProps := flattenApplicationInsightsAPIKeyLinkedProperties(result.LinkedReadProperties)
	if err := d.Set("read_permissions", readProps); err != nil {
		return fmt.Errorf("flattening `read_permissions `: %s", err)
	}
	writeProps := flattenApplicationInsightsAPIKeyLinkedProperties(result.LinkedWriteProperties)
	if err := d.Set("write_permissions", writeProps); err != nil {
		return fmt.Errorf("flattening `write_permissions `: %s", err)
	}

	return nil
}

func resourceApplicationInsightsAPIKeyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.APIKeysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiKeyID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting AzureRM Application Insights API key '%s'", id)

	result, err := client.Delete(ctx, id.ResourceGroup, id.ComponentName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(result.Response) {
			return nil
		}
		return fmt.Errorf("issuing AzureRM delete request for Application Insights API key '%s': %+v", id, err)
	}

	return nil
}

func camelCaseApiKeys(id string) string {
	// Azure only returns the api key identifier in the resource ID string where apikeys isn't camel cased
	r := regexp.MustCompile(`apikeys`)
	return r.ReplaceAllString(id, "apiKeys")
}
