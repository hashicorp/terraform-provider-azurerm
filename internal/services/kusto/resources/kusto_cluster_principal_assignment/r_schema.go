// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_cluster_principal_assignment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/clusterprincipalassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func kustoClusterPrincipalAssignmentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"cluster_name": {
			Description:  "The name of the cluster in which to create the resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ClusterName,
		},

		"name": {
			Description:  "The name of the Kusto Cluster Principal Assignment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ClusterPrincipalAssignmentName,
		},

		"tenant_id": {
			Description:  "The tenant ID in which the principal resides.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tenant_name": {
			Description: "The name of the tenant.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"principal_id": {
			Description:  "The object ID of the principal.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"principal_name": {
			Description: "The name of the principal.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"principal_type": {
			Description:  "The type of the principal.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(clusterprincipalassignments.PossibleValuesForPrincipalType(), false),
		},

		"role": {
			Description:  "The cluster role assigned to the principal.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(clusterprincipalassignments.PossibleValuesForClusterPrincipalRole(), false),
		},
	}
}
