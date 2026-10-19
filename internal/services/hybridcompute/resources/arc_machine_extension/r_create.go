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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MachineExtensionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.HybridCompute.MachineExtensionsClient
			machineId, err := machines.ParseMachineID(model.HybridComputeMachineId)
			if err != nil {
				return err
			}

			id := machineextensions.NewExtensionID(machineId.SubscriptionId, machineId.ResourceGroupName, machineId.MachineName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &machineextensions.MachineExtension{
				Location: location.Normalize(model.Location),
				Properties: &machineextensions.MachineExtensionProperties{
					EnableAutomaticUpgrade: &model.EnableAutomaticUpgrade,
				},
				Tags: &model.Tags,
			}

			if model.ForceUpdateTag != "" {
				properties.Properties.ForceUpdateTag = &model.ForceUpdateTag
			}

			if model.ProtectedSettings != "" {
				protectedSettingsValue := make(map[string]any)
				if err = json.Unmarshal([]byte(model.ProtectedSettings), &protectedSettingsValue); err != nil {
					return err
				}
				properties.Properties.ProtectedSettings = &protectedSettingsValue
			}

			if model.Publisher != "" {
				properties.Properties.Publisher = &model.Publisher
			}

			if model.Settings != "" {
				settingsValue := make(map[string]any)
				if err = json.Unmarshal([]byte(model.Settings), &settingsValue); err != nil {
					return err
				}
				properties.Properties.Settings = &settingsValue
			}

			if model.Type != "" {
				properties.Properties.Type = &model.Type
			}

			if model.TypeHandlerVersion != "" {
				properties.Properties.TypeHandlerVersion = &model.TypeHandlerVersion
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
