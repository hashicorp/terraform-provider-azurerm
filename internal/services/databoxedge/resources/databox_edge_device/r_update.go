// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databox_edge_device

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/databoxedge/2022-03-01/devices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataboxEdge.DeviceClient

			id, err := devices.ParseDataBoxEdgeDeviceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var metaModel EdgeDeviceModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := devices.DataBoxEdgeDevicePatch{}
			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = &metaModel.Tags
			}

			if _, err := client.Update(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
