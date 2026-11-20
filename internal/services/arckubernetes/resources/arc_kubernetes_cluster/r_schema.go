// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_cluster

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func arcKubernetesClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "Specifies the name which should be used for this Arc Kubernetes Cluster.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-_a-zA-Z0-9]{1,260}$"),
				"The name of Arc Kubernetes Cluster can only include alphanumeric characters, underscores, hyphens, has a maximum length of 260 characters, and must be unique.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"agent_public_key_certificate": {
			Description:  "Specifies the base64-encoded public certificate used by the agent to do the initial handshake to the backend services in Azure.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsBase64,
		},

		"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

		"location": commonschema.Location(),

		"agent_version": {
			Description: "Version of the agent running on the cluster resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"distribution": {
			Description: "The distribution running on this Arc Kubernetes Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"infrastructure": {
			Description: "The infrastructure on which the Arc Kubernetes Cluster is running on.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"kubernetes_version": {
			Description: "The Kubernetes version of the cluster resource.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"offering": {
			Description: "The cluster offering.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"total_core_count": {
			Description: "Number of CPU cores present in the cluster resource.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"total_node_count": {
			Description: "Number of nodes present in the cluster resource.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
