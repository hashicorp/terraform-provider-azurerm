// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_automatic_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.Containers.KubernetesClustersClient

			id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := clusterClient.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			if existing.Model.Identity != nil && existing.Model.Identity.IdentityIds != nil {
				for k := range existing.Model.Identity.IdentityIds {
					existing.Model.Identity.IdentityIds[k] = identity.UserAssignedIdentityDetails{}
				}
			}

			var model KubernetesAutomaticClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %w", err)
			}

			props := existing.Model.Properties

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = tags.Expand(model.Tags)
			}

			if metadata.ResourceData.HasChanges("api_server_access", "private_cluster") {
				props.ApiServerAccessProfile = expandKubernetesAutomaticClusterAPIAccessProfile(model)
			}

			if metadata.ResourceData.HasChange("identity") {
				existing.Model.Identity, err = identity.ExpandSystemOrUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding identity: %+v", err)
				}
			}

			if metadata.ResourceData.HasChange("service_mesh") {
				props.ServiceMeshProfile = expandKubernetesAutomaticClusterServiceMeshProfile(model.ServiceMeshProfile, props.ServiceMeshProfile)
			}

			if metadata.ResourceData.HasChange("web_app_routing_ingress") {
				props.IngressProfile = expandKubernetesAutomaticClusterWebAppRoutingIngress(model.WebAppRoutingIngress)
			}

			if err := clusterClient.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, managedclusters.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %w", *id, err)
			}

			return nil
		},
	}
}
