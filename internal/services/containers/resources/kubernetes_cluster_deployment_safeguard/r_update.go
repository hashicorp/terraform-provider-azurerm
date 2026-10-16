// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_deployment_safeguard

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-07-01/deploymentsafeguards"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.DeploymentSafeguardsClient

			kubernetesClusterId, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model KubernetesClusterDeploymentSafeguardResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			scopeId := commonids.NewScopeID(kubernetesClusterId.ID())

			existing, err := client.Get(ctx, scopeId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", kubernetesClusterId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", kubernetesClusterId)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", kubernetesClusterId)
			}

			payload := *existing.Model

			if metadata.ResourceData.HasChange("level") {
				payload.Properties.Level = deploymentsafeguards.DeploymentSafeguardsLevel(model.Level)
			}

			if metadata.ResourceData.HasChange("excluded_namespaces") {
				payload.Properties.ExcludedNamespaces = pointer.To(model.ExcludedNamespaces)
			}

			if metadata.ResourceData.HasChange("pod_security_standards_level") {
				payload.Properties.PodSecurityStandardsLevel = pointer.ToEnum[deploymentsafeguards.PodSecurityStandardsLevel](model.PodSecurityStandardsLevel)
			}

			if err := client.CreateThenPoll(ctx, scopeId, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", kubernetesClusterId, err)
			}

			return nil
		},
	}
}
