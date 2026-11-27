// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_blob

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_output_blob/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsOutputBlobCreateUpdate,
		Read:   resourceStreamAnalyticsOutputBlobRead,
		Update: resourceStreamAnalyticsOutputBlobCreateUpdate,
		Delete: resourceStreamAnalyticsOutputBlobDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := outputs.ParseOutputID(id)
			return err
		}, helpers.ImportStreamAnalyticsOutput(outputs.BlobOutputDataSource{})),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsOutputBlobV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: streamAnalyticsOutputBlobSchema(),
	}
}
