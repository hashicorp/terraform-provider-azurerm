// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_provisioned_cluster

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource             = Resource{}
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

// This resource is same type as the ArcKubernetesClusterResource but with kind="ProvisionedCluster".
type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_arc_kubernetes_provisioned_cluster"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return connectedclusters.ValidateConnectedClusterID
}
