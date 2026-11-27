// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_dapr_component

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/daprcomponents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	containerappsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/client"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := containerappsClient.From(metadata.Client).DaprComponentsClient
			id, err := daprcomponents.ParseDaprComponentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
