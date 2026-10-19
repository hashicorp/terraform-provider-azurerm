// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_web_test

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsWebTestsUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.WebTestsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := webtestsapis.ParseWebTestID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.WebTestsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	webTest := *existing.Model
	props := webTest.Properties

	if d.HasChange("description") {
		props.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("enabled") {
		props.Enabled = pointer.To(d.Get("enabled").(bool))
	}

	if d.HasChange("frequency") {
		props.Frequency = pointer.To(int64(d.Get("frequency").(int)))
	}

	if d.HasChange("timeout") {
		props.Timeout = pointer.To(int64(d.Get("timeout").(int)))
	}

	if d.HasChange("retry_enabled") {
		props.RetryEnabled = pointer.To(d.Get("retry_enabled").(bool))
	}

	if d.HasChange("geo_locations") {
		props.Locations = expandApplicationInsightsWebTestGeoLocations(d.Get("geo_locations").([]any))
	}

	if d.HasChange("configuration") {
		props.Configuration = &webtestsapis.WebTestPropertiesConfiguration{
			WebTest: pointer.To(d.Get("configuration").(string)),
		}
	}

	// Azure uses a special "hidden-link" tag to associate a web test with its parent Application Insights
	// component. The tag key is "hidden-link:<resource_id>" where <resource_id> is the full ARM resource ID
	// of the Application Insights component, and the value is "Resource". This tag is injected into the
	// user-supplied tags map before sending the request. It is generally undocumented but can be seen in
	// https://learn.microsoft.com/azure/azure-monitor/app/availability?tabs=standard
	if d.HasChange("tags") {
		appInsightsId, err := componentsapis.ParseComponentID(d.Get("application_insights_id").(string))
		if err != nil {
			return err
		}

		t := d.Get("tags").(map[string]any)
		tagKey := fmt.Sprintf("hidden-link:%s", appInsightsId.ID())
		t[tagKey] = "Resource"
		webTest.Tags = tags.Expand(t)
	}

	if _, err = client.WebTestsCreateOrUpdate(ctx, *id, webTest); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceApplicationInsightsWebTestsRead(d, meta)
}
