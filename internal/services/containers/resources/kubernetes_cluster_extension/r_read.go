// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_extension

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2024-11-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesExtensionsClient

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

			kubernetesClusterId, err := commonids.ParseKubernetesClusterID(id.Scope)
			if err != nil {
				return fmt.Errorf("parsing %q as a Kubernetes Cluster ID: %+v", id.Scope, err)
			}
			state := KubernetesClusterExtensionModel{
				Name:      id.ExtensionName,
				ClusterID: kubernetesClusterId.ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					var originalModel KubernetesClusterExtensionModel
					if err := metadata.Decode(&originalModel); err != nil {
						return fmt.Errorf("decoding: %+v", err)
					}

					if err = metadata.ResourceData.Set("aks_assigned_identity", flattenAksAssignedIdentity(properties.AksAssignedIdentity)); err != nil {
						return fmt.Errorf("setting `aks_assigned_identity`: %+v", err)
					}

					state.ConfigurationProtectedSettings = originalModel.ConfigurationProtectedSettings
					state.ConfigurationSettings = pointer.From(properties.ConfigurationSettings)
					state.CurrentVersion = pointer.From(properties.CurrentVersion)
					state.ExtensionType = pointer.From(properties.ExtensionType)
					state.Plan = flattenPlanModel(model.Plan)
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

func flattenPlanModel(input *extensions.Plan) []PlanModel {
	var outputList []PlanModel
	if input == nil {
		return []PlanModel{}
	}
	output := PlanModel{
		Name:      input.Name,
		Product:   input.Product,
		Publisher: input.Publisher,
	}

	output.PromotionCode = pointer.From(input.PromotionCode)
	output.Version = pointer.From(input.Version)

	return append(outputList, output)
}

func flattenAksAssignedIdentity(input *extensions.ExtensionPropertiesAksAssignedIdentity) []any {
	if input == nil {
		return []any{}
	}

	output := identity.SystemAssigned{
		Type:        identity.TypeSystemAssigned,
		PrincipalId: pointer.From(input.PrincipalId),
		TenantId:    pointer.From(input.TenantId),
	}

	if input.Type != nil && *input.Type == extensions.AKSIdentityTypeUserAssigned {
		output.Type = identity.TypeUserAssigned
	}

	return identity.FlattenSystemAssigned(&output)
}
