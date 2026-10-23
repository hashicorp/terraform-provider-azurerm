// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_machine_extension

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machineextensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridCompute.MachineExtensionsClient

			id, err := machineextensions.ParseExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model MachineExtensionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("automatic_upgrade_enabled") {
				properties.Properties.EnableAutomaticUpgrade = &model.EnableAutomaticUpgrade
			}

			if metadata.ResourceData.HasChange("force_update_tag") {
				if model.ForceUpdateTag != "" {
					properties.Properties.ForceUpdateTag = &model.ForceUpdateTag
				} else {
					properties.Properties.ForceUpdateTag = nil
				}
			}

			if metadata.ResourceData.HasChange("protected_settings") {
				protectedSettingsValue := make(map[string]any)
				if err := json.Unmarshal([]byte(model.ProtectedSettings), &protectedSettingsValue); err != nil {
					return err
				}

				properties.Properties.ProtectedSettings = &protectedSettingsValue
			}

			if metadata.ResourceData.HasChange("publisher") {
				if model.Publisher != "" {
					properties.Properties.Publisher = &model.Publisher
				} else {
					properties.Properties.Publisher = nil
				}
			}

			if metadata.ResourceData.HasChange("settings") {
				settingsValue := make(map[string]any)
				if err := json.Unmarshal([]byte(model.Settings), &settingsValue); err != nil {
					return err
				}

				properties.Properties.Settings = &settingsValue
			}

			if metadata.ResourceData.HasChange("type") {
				if model.Type != "" {
					properties.Properties.Type = &model.Type
				} else {
					properties.Properties.Type = nil
				}
			}

			if metadata.ResourceData.HasChange("type_handler_version") {
				if model.TypeHandlerVersion != "" {
					properties.Properties.TypeHandlerVersion = &model.TypeHandlerVersion
				} else {
					properties.Properties.TypeHandlerVersion = nil
				}
			}

			properties.SystemData = nil

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
