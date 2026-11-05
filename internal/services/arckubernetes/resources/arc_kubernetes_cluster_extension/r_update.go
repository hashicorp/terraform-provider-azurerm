// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_cluster_extension

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
			client := metadata.Client.ArcKubernetes.ExtensionsClient

			id, err := extensions.ParseScopedExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ArcKubernetesClusterExtensionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
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
