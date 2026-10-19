// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_device_update_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceupdate/2022-10-01/deviceupdates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.DeviceUpdatesClient

			id, err := deviceupdates.ParseInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model IotHubDeviceUpdateInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.InstancesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			existing := resp.Model
			if existing == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			// connectionString is not returned by API, so always expands DiagnosticStorageAccount
			expandedStorageProps, err := expandDiagnosticStorageAccount(metadata, model.DiagnosticStorageAccount)
			if err != nil {
				return fmt.Errorf("expanding `diagnostic_storage_account`: %w", err)
			}
			existing.Properties.DiagnosticStorageProperties = expandedStorageProps

			if metadata.ResourceData.HasChange("diagnostic_enabled") {
				existing.Properties.EnableDiagnostics = &model.DiagnosticEnabled
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = &model.Tags
			}

			if err := client.InstancesCreateThenPoll(ctx, *id, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
