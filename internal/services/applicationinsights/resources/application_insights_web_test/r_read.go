// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_web_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsWebTestsRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.WebTestsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webtestsapis.ParseWebTestID(d.Id())
	if err != nil {
		return err
	}

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

		parsedAppInsightsId, err := webtestsapis.ParseComponentIDInsensitively(appInsightsId)
		if err != nil {
			return err
		}
		d.Set("application_insights_id", parsedAppInsightsId.ID())

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenApplicationInsightsWebTestGeoLocations(input []webtestsapis.WebTestGeolocation) []string {
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
