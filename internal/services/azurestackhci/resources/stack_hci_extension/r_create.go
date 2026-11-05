// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_extension

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.Extensions

			var config StackHCIExtensionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			arcSettingId, err := extensions.ParseArcSettingID(config.ArcSettingId)
			if err != nil {
				return err
			}

			id := extensions.NewExtensionID(arcSettingId.SubscriptionId, arcSettingId.ResourceGroupName, arcSettingId.ClusterName, arcSettingId.ArcSettingName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			input := extensions.Extension{
				Properties: &extensions.ExtensionProperties{
					ExtensionParameters: &extensions.ExtensionParameters{
						AutoUpgradeMinorVersion: pointer.To(config.AutoUpgradeMinorVersionEnabled),
						EnableAutomaticUpgrade:  pointer.To(config.AutomaticUpgradeEnabled),
						Publisher:               pointer.To(config.Publisher),
						Type:                    pointer.To(config.Type),
					},
				},
			}

			if config.TypeHandlerVersion != "" {
				input.Properties.ExtensionParameters.TypeHandlerVersion = pointer.To(config.TypeHandlerVersion)
			}

			if config.Settings != "" {
				expandedSetting, err := pluginsdk.ExpandJsonFromString(config.Settings)
				if err != nil {
					return fmt.Errorf("expanding `setting`: %+v", err)
				}

				input.Properties.ExtensionParameters.Settings = pointer.To(any(expandedSetting))
			}

			if config.ProtectedSettings != "" {
				expandedSetting, err := pluginsdk.ExpandJsonFromString(config.ProtectedSettings)
				if err != nil {
					return fmt.Errorf("expanding `protected_settings`: %+v", err)
				}

				input.Properties.ExtensionParameters.ProtectedSettings = pointer.To(any(expandedSetting))
			}

			if err := client.CreateCallbackThenPoll(ctx, id, input, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
