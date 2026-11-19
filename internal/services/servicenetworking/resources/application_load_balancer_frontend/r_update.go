// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_frontend

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/frontendsinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (f Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.FrontendsInterface

			id, err := frontendsinterface.ParseFrontendID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config FrontendsModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			update := frontendsinterface.FrontendUpdate{}

			if metadata.ResourceData.HasChange("tags") {
				update.Tags = tags.Expand(config.Tags)
			}
			if _, err := client.Update(ctx, *id, update); err != nil {
				return fmt.Errorf("updating `azurerm_application_load_balancer_frontend` %s: %+v", *id, err)
			}

			return nil
		},
	}
}
