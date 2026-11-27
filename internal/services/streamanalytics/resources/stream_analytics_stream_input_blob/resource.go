// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_stream_input_blob

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_stream_input_blob/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsStreamInputBlobCreateUpdate,
		Read:   resourceStreamAnalyticsStreamInputBlobRead,
		Update: resourceStreamAnalyticsStreamInputBlobCreateUpdate,
		Delete: resourceStreamAnalyticsStreamInputBlobDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := inputs.ParseInputID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsStreamInputBlobV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: streamAnalyticsStreamInputBlobSchema(),
	}
}
