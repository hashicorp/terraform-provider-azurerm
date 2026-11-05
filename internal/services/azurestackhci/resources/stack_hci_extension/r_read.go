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

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.Extensions

			id, err := extensions.ParseExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			// protected_settings is not returned in the response, so we read it from the state
			var extension, config StackHCIExtensionResourceModel

			if err := metadata.Decode(&config); err != nil {
				return err
			}
			extension.ProtectedSettings = config.ProtectedSettings

			if model := existing.Model; model != nil {
				extension.Name = id.ExtensionName
				extension.ArcSettingId = extensions.NewArcSettingID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ArcSettingName).ID()

				if model.Properties != nil && model.Properties.ExtensionParameters != nil {
					param := model.Properties.ExtensionParameters
					extension.AutomaticUpgradeEnabled = pointer.From(param.EnableAutomaticUpgrade)
					extension.AutoUpgradeMinorVersionEnabled = pointer.From(param.AutoUpgradeMinorVersion)
					extension.Publisher = pointer.From(param.Publisher)
					extension.Type = pointer.From(param.Type)
					extension.TypeHandlerVersion = pointer.From(param.TypeHandlerVersion)

					var setting string
					if param.Settings != nil {
						setting, err = pluginsdk.FlattenJsonToString((*param.Settings).(map[string]any))
						if err != nil {
							return fmt.Errorf("flattening `settings`: %+v", err)
						}
					}
					extension.Settings = setting
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&extension)
		},
	}
}
