// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_time_series_database_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/timeseriesdatabaseconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := timeseriesdatabaseconnections.ParseTimeSeriesDatabaseConnectionID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.DigitalTwins.TimeSeriesDatabaseConnectionsClient
			if err = client.DeleteThenPoll(ctx, *id, timeseriesdatabaseconnections.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			return nil
		},
	}
}
