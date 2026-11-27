// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dashboard_grafana

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dashboard/2025-08-01/managedgrafanas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	dashboardClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/dashboard/client"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := dashboardClient.From(metadata.Client).GrafanaResourceClient

			id, err := managedgrafanas.ParseGrafanaID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.GrafanaDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
