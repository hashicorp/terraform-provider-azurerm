// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_deployment_safeguard

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-07-01/deploymentsafeguards"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"kubernetes_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKubernetesClusterID,
		},

		"level": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(deploymentsafeguards.PossibleValuesForDeploymentSafeguardsLevel(), false),
		},

		"excluded_namespaces": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"pod_security_standards_level": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      deploymentsafeguards.PodSecurityStandardsLevelPrivileged,
			ValidateFunc: validation.StringInSlice(deploymentsafeguards.PossibleValuesForPodSecurityStandardsLevel(), false),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
