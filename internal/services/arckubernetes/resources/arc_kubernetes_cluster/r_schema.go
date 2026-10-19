// Copyright IBM Corp. 2014, 2025
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-_a-zA-Z0-9]{1,260}$"),
				"The name of Arc Kubernetes Cluster can only include alphanumeric characters, underscores, hyphens, has a maximum length of 260 characters, and must be unique.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"agent_public_key_certificate": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsBase64,
		},

		"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

		"location": commonschema.Location(),

		"agent_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"distribution": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"infrastructure": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kubernetes_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"offering": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"total_core_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"total_node_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
