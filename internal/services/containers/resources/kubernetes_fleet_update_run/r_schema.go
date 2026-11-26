// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_update_run

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/updateruns"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for this Kubernetes Fleet Update Run.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"kubernetes_fleet_manager_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KubernetesFleetId{}, "Kubernetes Fleet Manager"),

		"managed_cluster_update": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"upgrade": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Description: "The type of upgrade to perform.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice([]string{
										string(updateruns.ManagedClusterUpgradeTypeFull),
										string(updateruns.ManagedClusterUpgradeTypeNodeImageOnly),
									}, false),
								},

								"kubernetes_version": {
									Description:  "The Kubernetes version to upgrade the member clusters to.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"node_image_selection": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Description: "The node image upgrade type.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice([]string{
										string(updateruns.NodeImageSelectionTypeConsistent),
										string(updateruns.NodeImageSelectionTypeLatest),
									}, false),
								},
							},
						},
					},
				},
			},
		},

		"fleet_update_strategy_id": {
			Description:   "The ID of the Fleet Update Strategy.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"stage"},
		},

		"stage": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ConflictsWith: []string{"fleet_update_strategy_id"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name which should be used for this stage.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"group": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description:  "The name which should be used for this group.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"after_stage_wait_in_seconds": {
						Description: "The time in seconds to wait at the end of this stage before starting the next one.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
