// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_insights_web_test

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func applicationInsightsWebTestSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Application Insights WebTest.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"application_insights_id": {
			Description:  "The ID of the Application Insights component on which the WebTest operates.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: componentsapis.ValidateComponentID,
		},

		"location": commonschema.Location(),

		"kind": {
			Description: "The kind of web test that this web test watches. Choices are `ping` and `multistep`.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(webtestsapis.WebTestKindMultistep),
				string(webtestsapis.WebTestKindPing),
			}, false),
		},

		"frequency": {
			Description: "Interval in seconds between test runs for this WebTest.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     300,
			ValidateFunc: validation.IntInSlice([]int{
				300,
				600,
				900,
			}),
		},

		"timeout": {
			Description: "Seconds until this WebTest will timeout and fail.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     30,
		},

		"enabled": {
			Description: "Is the test actively being monitored?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"retry_enabled": {
			Description: "Allow for retries should this WebTest fail.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"geo_locations": {
			Description: "A list of where to physically run the tests from to give global coverage for accessibility of your application.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:             pluginsdk.TypeString,
				ValidateFunc:     validation.StringIsNotEmpty,
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},
		},

		"description": {
			Description: "Purpose/user defined descriptive test for this WebTest.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"configuration": {
			Description:      "An XML configuration specification for a WebTest ([see here for more information](https://docs.microsoft.com/rest/api/application-insights/webtests/createorupdate/)).",
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: suppress.XmlDiff,
		},

		"tags": commonschema.Tags(),

		"synthetic_monitor_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
