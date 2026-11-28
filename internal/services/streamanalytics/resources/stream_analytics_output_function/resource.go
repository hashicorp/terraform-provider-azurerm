// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_function

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_output_function/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithCustomImporter = Resource{}
	_ sdk.ResourceWithStateMigration = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_stream_analytics_output_function"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return outputs.ValidateOutputID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.StreamAnalytics.OutputsClient
		resp, err := client.Get(ctx, *id)
		if err != nil || resp.Model == nil || resp.Model.Properties == nil {
			return fmt.Errorf("reading %s: %+v", *id, err)
		}

		props := resp.Model.Properties
		if _, ok := props.Datasource.(outputs.AzureFunctionOutputDataSource); !ok {
			return errors.New("specified output is not of type")
		}
		return nil
	}
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsOutputFunctionV0ToV1{},
		},
	}
}
