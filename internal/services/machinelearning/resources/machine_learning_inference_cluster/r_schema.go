// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_inference_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/machinelearningcomputes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func machineLearningInferenceClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"kubernetes_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKubernetesClusterID,
			// remove in 3.0 of the provider
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"location": commonschema.Location(),

		"machine_learning_workspace_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"cluster_purpose": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(machinelearningcomputes.ClusterPurposeFastProd),
			ValidateFunc: validation.StringInSlice(machinelearningcomputes.PossibleValuesForClusterPurpose(), false),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptionalForceNew(),

		"ssl": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cert": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ForceNew:      true,
						Default:       "",
						ConflictsWith: []string{"ssl.0.leaf_domain_label", "ssl.0.overwrite_existing_domain"},
					},
					"key": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ForceNew:      true,
						Default:       "",
						ConflictsWith: []string{"ssl.0.leaf_domain_label", "ssl.0.overwrite_existing_domain"},
					},
					"cname": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ForceNew:      true,
						Default:       "",
						ConflictsWith: []string{"ssl.0.leaf_domain_label", "ssl.0.overwrite_existing_domain"},
					},
					"leaf_domain_label": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ForceNew:      true,
						Default:       "",
						ConflictsWith: []string{"ssl.0.cert", "ssl.0.key", "ssl.0.cname"},
					},
					"overwrite_existing_domain": {
						Type:          pluginsdk.TypeBool,
						Optional:      true,
						ForceNew:      true,
						ConflictsWith: []string{"ssl.0.cert", "ssl.0.key", "ssl.0.cname"},
					},
				},
			},
		},

		"tags": commonschema.TagsForceNew(),
	}
}
