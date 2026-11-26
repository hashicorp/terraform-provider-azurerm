// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_machine_extension

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machineextensions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridCompute.MachineExtensionsClient

			id, err := machineextensions.ParseExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := MachineExtensionModel{
				Name:                   id.ExtensionName,
				HybridComputeMachineId: machines.NewMachineID(id.SubscriptionId, id.ResourceGroupName, id.MachineName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if properties := model.Properties; properties != nil {
					if properties.EnableAutomaticUpgrade != nil {
						state.EnableAutomaticUpgrade = *properties.EnableAutomaticUpgrade
					}

					if properties.ForceUpdateTag != nil {
						state.ForceUpdateTag = *properties.ForceUpdateTag
					}

					if properties.Publisher != nil {
						state.Publisher = *properties.Publisher
					}

					var extModel MachineExtensionModel
					if err := metadata.Decode(&extModel); err != nil {
						return err
					}

					if extModel.ProtectedSettings != "" {
						state.ProtectedSettings = extModel.ProtectedSettings
					}

					if properties.Settings != nil && *properties.Settings != nil {
						settingsValue, err := json.Marshal(*properties.Settings)
						if err != nil {
							return err
						}

						state.Settings = string(settingsValue)
					}

					if properties.Type != nil {
						state.Type = *properties.Type
					}

					if properties.TypeHandlerVersion != nil {
						state.TypeHandlerVersion = *properties.TypeHandlerVersion
					}
				}

				if model.Tags != nil {
					state.Tags = *model.Tags
				}
			}

			return metadata.Encode(&state)
		},
	}
}
