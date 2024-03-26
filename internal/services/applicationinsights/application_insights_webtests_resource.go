// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	webtests "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsWebTests() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApplicationInsightsWebTestsCreateUpdate,
		Read:   resourceApplicationInsightsWebTestsRead,
		Update: resourceApplicationInsightsWebTestsCreateUpdate,
		Delete: resourceApplicationInsightsWebTestsDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := webtests.ParseWebTestID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WebTestUpgradeV0ToV1{},
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

			"resource_group_name": commonschema.ResourceGroupName(),

			"application_insights_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: components.ValidateComponentID,
			},

			"location": commonschema.Location(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(webtests.WebTestKindMultistep),
					string(webtests.WebTestKindPing),
				}, false),
			},

			"frequency": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  300,
				ValidateFunc: validation.IntInSlice([]int{
					300,
					600,
					900,
				}),
			},

			"timeout": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  30,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"retry_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"geo_locations": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:             pluginsdk.TypeString,
					ValidateFunc:     validation.StringIsNotEmpty,
					StateFunc:        location.StateFunc,
					DiffSuppressFunc: location.DiffSuppressFunc,
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"configuration": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.XmlDiff,
			},

			"tags": commonschema.Tags(),

			"synthetic_monitor_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApplicationInsightsWebTestsCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.WebTestsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Application Insights WebTest creation.")

	appInsightsId, err := components.ParseComponentID(d.Get("application_insights_id").(string))
	if err != nil {
		return err
	}

	id := webtests.NewWebTestID(appInsightsId.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.WebTestsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_application_insights_web_test", id.ID())
		}
	}

	kind := d.Get("kind").(string)
	description := d.Get("description").(string)
	frequency := int32(d.Get("frequency").(int))
	timeout := int32(d.Get("timeout").(int))
	isEnabled := d.Get("enabled").(bool)
	retryEnabled := d.Get("retry_enabled").(bool)
	geoLocationsRaw := d.Get("geo_locations").([]interface{})
	geoLocations := expandApplicationInsightsWebTestGeoLocations(geoLocationsRaw)
	testConf := d.Get("configuration").(string)

	t := d.Get("tags").(map[string]interface{})
	tagKey := fmt.Sprintf("hidden-link:%s", appInsightsId.ID())
	t[tagKey] = "Resource"

	webTest := webtests.WebTest{
		Name:     pointer.To(id.WebTestName),
		Location: location.Normalize(d.Get("location").(string)),
		Kind:     pointer.To(webtests.WebTestKind(kind)),
		Properties: &webtests.WebTestProperties{
			SyntheticMonitorId: id.WebTestName,
			Name:               id.WebTestName,
			Description:        &description,
			Enabled:            &isEnabled,
			Frequency:          pointer.To(int64(frequency)),
			Timeout:            pointer.To(int64(timeout)),
			Kind:               webtests.WebTestKind(kind),
			RetryEnabled:       &retryEnabled,
			Locations:          geoLocations,
			Configuration: &webtests.WebTestPropertiesConfiguration{
				WebTest: &testConf,
			},
		},
		Tags: tags.Expand(t),
	}

	_, err = client.WebTestsCreateOrUpdate(ctx, id, webTest)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApplicationInsightsWebTestsRead(d, meta)
}

func resourceApplicationInsightsWebTestsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.WebTestsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webtests.ParseWebTestID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading AzureRM Application Insights %q", *id)

	resp, err := client.WebTestsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.WebTestName)
	d.Set("resource_group_name", id.ResourceGroupName)

	appInsightsId := ""

	if model := resp.Model; model != nil {
		if model.Tags != nil {
			for i := range *model.Tags {
				if strings.HasPrefix(i, "hidden-link") {
					appInsightsId = strings.Split(i, ":")[1]
				}
			}
		}
		d.Set("kind", pointer.From(model.Kind))
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			// It is possible that the root level `kind` in response is empty in some cases (see PR #8372 for more info)
			if model.Kind == nil || *model.Kind == "" {
				d.Set("kind", props.Kind)
			}
			d.Set("synthetic_monitor_id", props.SyntheticMonitorId)
			d.Set("description", props.Description)
			d.Set("enabled", props.Enabled)
			d.Set("frequency", props.Frequency)
			d.Set("timeout", props.Timeout)
			d.Set("retry_enabled", props.RetryEnabled)

			if config := props.Configuration; config != nil {
				d.Set("configuration", config.WebTest)
			}

			if err := d.Set("geo_locations", flattenApplicationInsightsWebTestGeoLocations(props.Locations)); err != nil {
				return fmt.Errorf("setting `geo_locations`: %+v", err)
			}
		}

		parsedAppInsightsId, err := webtests.ParseComponentIDInsensitively(appInsightsId)
		if err != nil {
			return err
		}
		d.Set("application_insights_id", parsedAppInsightsId.ID())

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceApplicationInsightsWebTestsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.WebTestsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webtests.ParseWebTestID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.WebTestsDelete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return err
}

func expandApplicationInsightsWebTestGeoLocations(input []interface{}) []webtests.WebTestGeolocation {
	locations := make([]webtests.WebTestGeolocation, 0)

	for _, v := range input {
		lc := v.(string)
		loc := webtests.WebTestGeolocation{
			Id: &lc,
		}
		locations = append(locations, loc)
	}

	return locations
}

func flattenApplicationInsightsWebTestGeoLocations(input []webtests.WebTestGeolocation) []string {
	results := make([]string, 0)
	if len(input) == 0 {
		return results
	}

	for _, prop := range input {
		if prop.Id != nil {
			results = append(results, location.Normalize(*prop.Id))
		}
	}

	return results
}
