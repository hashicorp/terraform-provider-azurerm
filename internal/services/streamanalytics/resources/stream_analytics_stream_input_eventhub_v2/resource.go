// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_stream_input_eventhub_v2

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_stream_input_eventhub_v2/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithCustomImporter = Resource{}
	_ sdk.ResourceWithStateMigration = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_stream_analytics_stream_input_eventhub_v2"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return inputs.ValidateInputID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := inputs.ParseInputID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.StreamAnalytics.InputsClient
		resp, err := client.Get(ctx, *id)
		if err != nil || resp.Model == nil || resp.Model.Properties == nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		props := resp.Model.Properties

		streamInput, ok := props.(inputs.StreamInputProperties)
		if !ok {
			return fmt.Errorf("specified resource is not a Stream Input: %+v", err)
		}

		if _, ok := streamInput.Datasource.(inputs.EventHubV2StreamInputDataSource); !ok {
			return fmt.Errorf("specified input is not of type EventHubV2: %+v", err)
		}

		return nil
	}
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsStreamInputEventHubV2V0ToV1{},
		},
	}
}
