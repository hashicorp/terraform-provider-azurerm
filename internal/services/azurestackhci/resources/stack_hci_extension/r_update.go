// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_extension

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.Extensions

			var config StackHCIExtensionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			id, err := extensions.ParseExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model

			if model == nil || model.Properties == nil || model.Properties.ExtensionParameters == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			updateModel := extensions.ExtensionPatch{
				Properties: &extensions.ExtensionPatchProperties{
					ExtensionParameters: &extensions.ExtensionPatchParameters{},
				},
			}

			if metadata.ResourceData.HasChange("automatic_upgrade_enabled") {
				updateModel.Properties.ExtensionParameters.EnableAutomaticUpgrade = pointer.To(config.AutomaticUpgradeEnabled)
			}

			if metadata.ResourceData.HasChange("protected_settings") {
				if config.ProtectedSettings != "" {
					expandedSetting, err := pluginsdk.ExpandJsonFromString(config.ProtectedSettings)
					if err != nil {
						return fmt.Errorf("expanding `protected_settings`: %+v", err)
					}

					updateModel.Properties.ExtensionParameters.ProtectedSettings = pointer.To(any(expandedSetting))
				} else {
					var emptyInterface any
					updateModel.Properties.ExtensionParameters.Settings = pointer.To(emptyInterface)
				}
			}

			if metadata.ResourceData.HasChange("settings") {
				if config.Settings != "" {
					expandedSetting, err := pluginsdk.ExpandJsonFromString(config.Settings)
					if err != nil {
						return fmt.Errorf("expanding `setting`: %+v", err)
					}

					updateModel.Properties.ExtensionParameters.Settings = pointer.To(any(expandedSetting))
				} else {
					var emptyInterface any
					updateModel.Properties.ExtensionParameters.Settings = pointer.To(emptyInterface)
				}
			}

			if metadata.ResourceData.HasChange("type_handler_version") {
				updateModel.Properties.ExtensionParameters.TypeHandlerVersion = pointer.To(config.TypeHandlerVersion)
			}

			if err := client.UpdateThenPoll(ctx, *id, updateModel); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
