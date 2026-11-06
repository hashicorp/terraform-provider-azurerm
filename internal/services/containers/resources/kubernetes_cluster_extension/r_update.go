// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_extension

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2024-11-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesExtensionsClient

			id, err := extensions.ParseScopedExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model KubernetesClusterExtensionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("configuration_settings") {
				previous, _ := metadata.ResourceData.GetChange("configuration_settings")
				if kubernetesClusterExtensionSettingsRemoved(previous.(map[string]any), model.ConfigurationSettings) {
					existing, err := client.Get(ctx, *id)
					if err != nil {
						return fmt.Errorf("retrieving %s: %+v", *id, err)
					}
					if existing.Model == nil {
						return fmt.Errorf("retrieving %s: missing extension", *id)
					}

					previousProtected, _ := metadata.ResourceData.GetChange("configuration_protected_settings")
					properties, err := expandKubernetesClusterExtensionFullUpdate(*existing.Model, model.ConfigurationSettings, model.ConfigurationProtectedSettings, previousProtected.(map[string]any))
					if err != nil {
						return err
					}

					// PATCH merges the map, so omitting a key does not remove it. PUT
					// replaces it while retaining the existing extension's lifecycle.
					if err := client.CreateThenPoll(ctx, *id, properties); err != nil {
						return fmt.Errorf("updating %s: %+v", *id, err)
					}
					return nil
				}
			}

			properties := &extensions.PatchExtension{
				Properties: &extensions.PatchExtensionProperties{},
			}

			if metadata.ResourceData.HasChange("configuration_protected_settings") {
				properties.Properties.ConfigurationProtectedSettings = &model.ConfigurationProtectedSettings
			}

			if metadata.ResourceData.HasChange("configuration_settings") {
				properties.Properties.ConfigurationSettings = &model.ConfigurationSettings
			}

			if err := client.UpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
