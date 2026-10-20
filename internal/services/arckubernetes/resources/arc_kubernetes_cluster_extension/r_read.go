// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_cluster_extension

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2024-11-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ExtensionsClient

			id, err := extensions.ParseScopedExtensionID(metadata.ResourceData.Id())
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

			clusterId, err := connectedclusters.ParseConnectedClusterID(id.Scope)
			if err != nil {
				return fmt.Errorf("parsing %q as a Connected Cluster ID: %+v", id.Scope, err)
			}

			state := ArcKubernetesClusterExtensionModel{
				Name:      id.ExtensionName,
				ClusterID: clusterId.ID(),
			}

			if model := resp.Model; model != nil {
				if err = metadata.ResourceData.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}

				if properties := model.Properties; properties != nil {
					var originalModel ArcKubernetesClusterExtensionModel
					if err := metadata.Decode(&originalModel); err != nil {
						return fmt.Errorf("decoding: %+v", err)
					}

					state.ConfigurationProtectedSettings = originalModel.ConfigurationProtectedSettings
					state.ConfigurationSettings = pointer.From(properties.ConfigurationSettings)
					state.CurrentVersion = pointer.From(properties.CurrentVersion)
					state.ExtensionType = pointer.From(properties.ExtensionType)
					state.ReleaseTrain = pointer.From(properties.ReleaseTrain)

					if properties.Scope != nil {
						if properties.Scope.Cluster != nil {
							state.ReleaseNamespace = pointer.From(properties.Scope.Cluster.ReleaseNamespace)
						}

						if properties.Scope.Namespace != nil {
							state.TargetNamespace = pointer.From(properties.Scope.Namespace.TargetNamespace)
						}
					}

					state.Version = pointer.From(properties.Version)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
