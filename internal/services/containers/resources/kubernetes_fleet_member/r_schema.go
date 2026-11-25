// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_member

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"kubernetes_cluster_id": {
			Description: "The ARM resource ID of the cluster that joins the Fleet.",
			ForceNew:    true,
			Required:    true,
			Type:        pluginsdk.TypeString,
		},
		"kubernetes_fleet_id": {
			Description: "The Kubernetes Fleet ID within which this Kubernetes Fleet Member should exist.",
			ForceNew:    true,
			Required:    true,
			Type:        pluginsdk.TypeString,
		},
		"name": {
			Description: "The name of this Kubernetes Fleet Member.",
			ForceNew:    true,
			Required:    true,
			Type:        pluginsdk.TypeString,
		},
		"group": {
			Description: "The group this member belongs to for multi-cluster update management.",
			Optional:    true,
			Type:        pluginsdk.TypeString,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
