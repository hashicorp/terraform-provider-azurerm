// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_deployment_safeguard

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-07-01/deploymentsafeguards"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.DeploymentSafeguardsClient

			var model KubernetesClusterDeploymentSafeguardResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			kubernetesClusterId, err := commonids.ParseKubernetesClusterID(model.KubernetesClusterId)
			if err != nil {
				return err
			}

			scopeId := commonids.NewScopeID(kubernetesClusterId.ID())

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, scopeId)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", kubernetesClusterId, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), kubernetesClusterId)
				}
			}

			payload := deploymentsafeguards.DeploymentSafeguard{
				Properties: &deploymentsafeguards.DeploymentSafeguardsProperties{
					Level:                     deploymentsafeguards.DeploymentSafeguardsLevel(model.Level),
					PodSecurityStandardsLevel: pointer.ToEnum[deploymentsafeguards.PodSecurityStandardsLevel](model.PodSecurityStandardsLevel),
				},
			}

			if len(model.ExcludedNamespaces) > 0 {
				payload.Properties.ExcludedNamespaces = pointer.To(model.ExcludedNamespaces)
			}

			if err := client.CreateCallbackThenPoll(ctx, scopeId, payload, metadata.SetIDCallback(kubernetesClusterId)); err != nil {
				return fmt.Errorf("creating %s: %+v", kubernetesClusterId, err)
			}

			metadata.SetID(kubernetesClusterId)
			return nil
		},
	}
}
