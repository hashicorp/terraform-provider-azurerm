// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2020-02-02/insights" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApplicationInsightsWebTests() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApplicationInsightsWebTestsCreateUpdate,
		Read:   resourceApplicationInsightsWebTestsRead,
		Update: resourceApplicationInsightsWebTestsCreateUpdate,
		Delete: resourceApplicationInsightsWebTestsDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WebTestID(id)
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
				ValidateFunc: validate.ComponentID,
			},

			"location": commonschema.Location(),

			"kind": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(insights.WebTestKindMultistep),
					string(insights.WebTestKindPing),
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

			"tags": tags.Schema(),

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

	appInsightsId, err := parse.ComponentID(d.Get("application_insights_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewWebTestID(appInsightsId.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Application Insights %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_application_insights_web_test", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
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

	webTest := insights.WebTest{
		Name:     &id.Name,
		Location: &location,
		Kind:     insights.WebTestKind(kind),
		WebTestProperties: &insights.WebTestProperties{
			SyntheticMonitorID: &id.Name,
			WebTestName:        &id.Name,
			Description:        &description,
			Enabled:            &isEnabled,
			Frequency:          &frequency,
			Timeout:            &timeout,
			WebTestKind:        insights.WebTestKind(kind),
			RetryEnabled:       &retryEnabled,
			Locations:          &geoLocations,
			Configuration: &insights.WebTestPropertiesConfiguration{
				WebTest: &testConf,
			},
		},
		Tags: tags.Expand(t),
	}

	_, err = client.CreateOrUpdate(ctx, id, webTest)
	if err != nil {
		return fmt.Errorf("creating/updating Application Insights %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApplicationInsightsWebTestsRead(d, meta)
}

func resourceApplicationInsightsWebTestsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.WebTestsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebTestID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading AzureRM Application Insights %q", *id)

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Application Insights %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Application Insights %s: %+v", *id, err)
	}

	appInsightsId := ""
	for i := range resp.Tags {
		if strings.HasPrefix(i, "hidden-link") {
			appInsightsId = strings.Split(i, ":")[1]
		}
	}
	parsedAppInsightsId, err := parse.ComponentIDInsensitively(appInsightsId)
	if err != nil {
		return fmt.Errorf("parsing `application_insights_id`: %+v", err)
	}

	d.Set("application_insights_id", parsedAppInsightsId.ID())
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("kind", resp.Kind)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.WebTestProperties; props != nil {
		// It is possible that the root level `kind` in response is empty in some cases (see PR #8372 for more info)
		if resp.Kind == "" {
			d.Set("kind", props.WebTestKind)
		}
		d.Set("synthetic_monitor_id", props.SyntheticMonitorID)
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

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceApplicationInsightsWebTestsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppInsights.WebTestsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WebTestID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting AzureRM Application Insights %s", *id)

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("issuing AzureRM delete request for Application Insights '%s': %+v", *id, err)
	}

	return err
}

func expandApplicationInsightsWebTestGeoLocations(input []interface{}) []insights.WebTestGeolocation {
	locations := make([]insights.WebTestGeolocation, 0)

	for _, v := range input {
		lc := v.(string)
		loc := insights.WebTestGeolocation{
			Location: &lc,
		}
		locations = append(locations, loc)
	}

	return locations
}

func flattenApplicationInsightsWebTestGeoLocations(input *[]insights.WebTestGeolocation) []string {
	results := make([]string, 0)
	if input == nil {
		return results
	}

	for _, prop := range *input {
		if prop.Location != nil {
			results = append(results, azure.NormalizeLocation(*prop.Location))
		}
	}

	return results
}
