// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_trusted_access_role_binding

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/trustedaccess"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.TrustedAccessClient
			schema := KubernetesClusterTrustedAccessRoleBindingResourceSchema{}

			id, err := trustedaccess.ParseTrustedAccessRoleBindingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			kubernetesClusterId := commonids.NewKubernetesClusterID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName)

			resp, err := client.RoleBindingsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.KubernetesClusterId = kubernetesClusterId.ID()
				schema.Name = id.TrustedAccessRoleBindingName
				r.mapTrustedAccessRoleBindingToKubernetesClusterTrustedAccessRoleBindingResourceSchema(*model, &schema)
			}

			return metadata.Encode(&schema)
		},
	}
}

type KubernetesClusterTrustedAccessRoleBindingResourceSchema struct {
	KubernetesClusterId string   `tfschema:"kubernetes_cluster_id"`
	Name                string   `tfschema:"name"`
	Roles               []string `tfschema:"roles"`
	SourceResourceId    string   `tfschema:"source_resource_id"`
}

func (r Resource) mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBinding(input KubernetesClusterTrustedAccessRoleBindingResourceSchema, output *trustedaccess.TrustedAccessRoleBinding) {
	r.mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBindingProperties(input, &output.Properties)
}

func (r Resource) mapTrustedAccessRoleBindingToKubernetesClusterTrustedAccessRoleBindingResourceSchema(input trustedaccess.TrustedAccessRoleBinding, output *KubernetesClusterTrustedAccessRoleBindingResourceSchema) {
	r.mapTrustedAccessRoleBindingPropertiesToKubernetesClusterTrustedAccessRoleBindingResourceSchema(input.Properties, output)
}

func (r Resource) mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBindingProperties(input KubernetesClusterTrustedAccessRoleBindingResourceSchema, output *trustedaccess.TrustedAccessRoleBindingProperties) {
	output.Roles = append([]string{}, input.Roles...)

	output.SourceResourceId = input.SourceResourceId
}

func (r Resource) mapTrustedAccessRoleBindingPropertiesToKubernetesClusterTrustedAccessRoleBindingResourceSchema(input trustedaccess.TrustedAccessRoleBindingProperties, output *KubernetesClusterTrustedAccessRoleBindingResourceSchema) {
	output.Roles = append([]string{}, input.Roles...)

	output.SourceResourceId = input.SourceResourceId
}
