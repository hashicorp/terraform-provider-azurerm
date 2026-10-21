// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_security_perimeter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecurityperimeters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.NetworkSecurityPerimetersClient

			id, err := networksecurityperimeters.ParseNetworkSecurityPerimeterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config NetworkSecurityPerimeterResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("tags") {
				tagReq := networksecurityperimeters.UpdateTagsRequest{
					Id:   pointer.To(id.ID()),
					Tags: pointer.To(config.Tags),
				}

				if _, err := client.Patch(ctx, *id, tagReq); err != nil {
					return fmt.Errorf("updating tags for %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
