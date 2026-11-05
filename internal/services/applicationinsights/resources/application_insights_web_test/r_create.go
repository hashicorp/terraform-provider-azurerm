// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_web_test

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApplicationInsightsWebTestsCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppInsights.WebTestsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appInsightsId, err := componentsapis.ParseComponentID(d.Get("application_insights_id").(string))
	if err != nil {
		return err
	}

	id := webtestsapis.NewWebTestID(appInsightsId.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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

	// Azure uses a special "hidden-link" tag to associate a web test with its parent Application Insights
	// component. The tag key is "hidden-link:<resource_id>" where <resource_id> is the full ARM resource ID
	// of the Application Insights component, and the value is "Resource". This tag is injected into the
	// user-supplied tags map before sending the request. It is generally undocumented but can be seen in
	// https://learn.microsoft.com/azure/azure-monitor/app/availability?tabs=standard
	t := d.Get("tags").(map[string]any)
	tagKey := fmt.Sprintf("hidden-link:%s", appInsightsId.ID())
	t[tagKey] = "Resource"

	webTest := webtestsapis.WebTest{
		Name:     pointer.To(id.WebTestName),
		Location: location.Normalize(d.Get("location").(string)),
		Kind:     pointer.ToEnum[webtestsapis.WebTestKind](d.Get("kind").(string)),
		Properties: &webtestsapis.WebTestProperties{
			SyntheticMonitorId: id.WebTestName,
			Name:               id.WebTestName,
			Description:        pointer.To(d.Get("description").(string)),
			Enabled:            pointer.To(d.Get("enabled").(bool)),
			Frequency:          pointer.To(int64(d.Get("frequency").(int))),
			Timeout:            pointer.To(int64(d.Get("timeout").(int))),
			Kind:               webtestsapis.WebTestKind(d.Get("kind").(string)),
			RetryEnabled:       pointer.To(d.Get("retry_enabled").(bool)),
			Locations:          expandApplicationInsightsWebTestGeoLocations(d.Get("geo_locations").([]any)),
			Configuration: &webtestsapis.WebTestPropertiesConfiguration{
				WebTest: pointer.To(d.Get("configuration").(string)),
			},
		},
		Tags: tags.Expand(t),
	}

	if _, err = client.WebTestsCreateOrUpdate(ctx, id, webTest); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceApplicationInsightsWebTestsRead(d, meta)
}

func expandApplicationInsightsWebTestGeoLocations(input []any) []webtestsapis.WebTestGeolocation {
	locations := make([]webtestsapis.WebTestGeolocation, 0)

	for _, v := range input {
		loc := webtestsapis.WebTestGeolocation{
			Id: pointer.To(v.(string)),
		}
		locations = append(locations, loc)
	}

	return locations
}
