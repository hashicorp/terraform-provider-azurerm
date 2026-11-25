// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster_trusted_access_role_binding

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"kubernetes_cluster_id": {
			Description: "The Kubernetes Cluster ID within which this Kubernetes Cluster Trusted Access Role Binding should exist.",
			ForceNew:    true,
			Required:    true,
			Type:        pluginsdk.TypeString,
		},
		"name": {
			Description: "The name of this Kubernetes Cluster Trusted Access Role Binding.",
			ForceNew:    true,
			Required:    true,
			Type:        pluginsdk.TypeString,
		},
		"roles": {
			Description: "A list of roles to bind, each item is a resource type qualified role name.",
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Required: true,
			Type:     pluginsdk.TypeList,
		},
		"source_resource_id": {
			Description: "The ARM resource ID of source resource that trusted access is configured for.",
			ForceNew:    true,
			Required:    true,
			Type:        pluginsdk.TypeString,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
