// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package graph_services_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/graphservices/2023-04-13/graphservicesprods"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Graph.V20230413.Graphservicesprods

			id, err := graphservicesprods.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.AccountsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
