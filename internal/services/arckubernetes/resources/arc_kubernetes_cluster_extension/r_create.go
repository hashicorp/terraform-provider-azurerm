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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ArcKubernetesClusterExtensionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.ArcKubernetes.ExtensionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			clusterID, err := connectedclusters.ParseConnectedClusterID(model.ClusterID)
			if err != nil {
				return err
			}

			// defined as strings because they're not enums in the swagger https://github.com/Azure/azure-rest-api-specs/pull/23545
			connectedClusterId := connectedclusters.NewConnectedClusterID(subscriptionId, clusterID.ResourceGroupName, clusterID.ConnectedClusterName)
			id := extensions.NewScopedExtensionID(connectedClusterId.ID(), model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &extensions.Extension{
				Properties: &extensions.ExtensionProperties{
					AutoUpgradeMinorVersion:        pointer.To(model.Version == ""),
					ConfigurationProtectedSettings: &model.ConfigurationProtectedSettings,
					ConfigurationSettings:          &model.ConfigurationSettings,
				},
			}

			identityValue, err := identity.ExpandSystemAssigned(metadata.ResourceData.Get("identity").([]any))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			properties.Identity = identityValue

			if model.ExtensionType != "" {
				properties.Properties.ExtensionType = &model.ExtensionType
			}

			if model.ReleaseNamespace != "" {
				properties.Properties.Scope = &extensions.Scope{
					Cluster: &extensions.ScopeCluster{
						ReleaseNamespace: &model.ReleaseNamespace,
					},
				}
			}

			if model.ReleaseTrain != "" {
				properties.Properties.ReleaseTrain = &model.ReleaseTrain
			}

			if model.TargetNamespace != "" {
				properties.Properties.Scope = &extensions.Scope{
					Namespace: &extensions.ScopeNamespace{
						TargetNamespace: &model.TargetNamespace,
					},
				}
			}

			if model.Version != "" {
				properties.Properties.Version = &model.Version
			}

			if err := client.CreateCallbackThenPoll(ctx, id, *properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
