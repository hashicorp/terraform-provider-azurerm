// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_ipam_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipampools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.IPamPools

			id, err := ipampools.ParseIPamPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerIpamPoolResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := ipampools.IPamPoolUpdate{
				Properties: &ipampools.IPamPoolUpdateProperties{},
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("description") {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("display_name") {
				parameters.Properties.DisplayName = pointer.To(model.DisplayName)
			}

			if _, err := client.Update(ctx, *id, parameters, ipampools.DefaultUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}
