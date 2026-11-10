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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.DeploymentSafeguardsClient

			kubernetesClusterId, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			scopeId := commonids.NewScopeID(kubernetesClusterId.ID())

			resp, err := client.Get(ctx, scopeId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(kubernetesClusterId)
				}
				return fmt.Errorf("retrieving %s: %+v", kubernetesClusterId, err)
			}

			state := KubernetesClusterDeploymentSafeguardResourceModel{
				KubernetesClusterId: kubernetesClusterId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Level = string(props.Level)
					state.ExcludedNamespaces = pointer.From(props.ExcludedNamespaces)
					state.PodSecurityStandardsLevel = pointer.FromEnum(props.PodSecurityStandardsLevel)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
