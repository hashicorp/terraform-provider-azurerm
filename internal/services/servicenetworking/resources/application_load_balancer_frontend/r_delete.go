// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_frontend

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/frontendsinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	servicenetworkingClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/client"
)

func (f Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := servicenetworkingClient.From(metadata.Client).FrontendsInterface

			id, err := frontendsinterface.ParseFrontendID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %q: %+v", id.ID(), err)
			}

			return nil
		},
	}
}
