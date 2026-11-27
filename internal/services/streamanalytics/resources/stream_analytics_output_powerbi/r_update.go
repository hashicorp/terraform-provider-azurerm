// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_powerbi

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.OutputsClient
			id, err := outputs.ParseOutputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state OutputPowerBIResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			needUpdateDataSourceProps := false
			dataSourceProps := outputs.PowerBIOutputDataSourceProperties{}
			d := metadata.ResourceData

			if d.HasChange("dataset") {
				needUpdateDataSourceProps = true
				dataSourceProps.Dataset = &state.DataSet
			}

			if d.HasChange("table") {
				needUpdateDataSourceProps = true
				dataSourceProps.Table = &state.Table
			}

			if d.HasChange("group_name") {
				needUpdateDataSourceProps = true
				dataSourceProps.GroupName = &state.GroupName
			}

			if d.HasChange("group_id") {
				needUpdateDataSourceProps = true
				dataSourceProps.GroupId = &state.GroupID
			}

			if d.HasChange("token_user_principal_name") {
				needUpdateDataSourceProps = true
				dataSourceProps.TokenUserPrincipalName = &state.TokenUserPrincipalName
			}

			if d.HasChange("token_user_display_name") {
				needUpdateDataSourceProps = true
				dataSourceProps.TokenUserDisplayName = &state.TokenUserDisplayName
			}

			if !needUpdateDataSourceProps {
				return nil
			}

			updateDataSource := outputs.PowerBIOutputDataSource{
				Properties: &dataSourceProps,
			}

			props := outputs.Output{
				Properties: &outputs.OutputProperties{
					Datasource: updateDataSource,
				},
			}

			var opts outputs.UpdateOperationOptions
			if _, err = client.Update(ctx, *id, props, opts); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
