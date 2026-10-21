// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_trusted_access_role_binding

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/trustedaccess"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.TrustedAccessClient

			var config KubernetesClusterTrustedAccessRoleBindingResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			kubernetesClusterId, err := commonids.ParseKubernetesClusterID(config.KubernetesClusterId)
			if err != nil {
				return err
			}

			id := trustedaccess.NewTrustedAccessRoleBindingID(subscriptionId, kubernetesClusterId.ResourceGroupName, kubernetesClusterId.ManagedClusterName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.RoleBindingsGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			var payload trustedaccess.TrustedAccessRoleBinding
			r.mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBinding(config, &payload)

			if err := client.RoleBindingsCreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
