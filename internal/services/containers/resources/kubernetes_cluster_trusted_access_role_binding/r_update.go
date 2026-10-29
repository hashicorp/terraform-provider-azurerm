// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_trusted_access_role_binding

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/trustedaccess"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.TrustedAccessClient

			id, err := trustedaccess.ParseTrustedAccessRoleBindingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config KubernetesClusterTrustedAccessRoleBindingResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.RoleBindingsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: properties was nil", *id)
			}
			payload := *existing.Model

			r.mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBinding(config, &payload)

			if _, err := client.RoleBindingsCreateOrUpdate(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
