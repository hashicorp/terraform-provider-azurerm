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
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"kubernetes_fleet_manager_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KubernetesFleetId{}),

		"managed_cluster_update": {
			Required: true,
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"upgrade": {
						Required: true,
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Required: true,
									Type:     pluginsdk.TypeString,
									ValidateFunc: validation.StringInSlice([]string{
										string(updateruns.ManagedClusterUpgradeTypeFull),
										string(updateruns.ManagedClusterUpgradeTypeNodeImageOnly),
									}, false),
								},

								"kubernetes_version": {
									Optional:     true,
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"node_image_selection": {
						Optional: true,
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Required: true,
									Type:     pluginsdk.TypeString,
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
			Optional:      true,
			Type:          pluginsdk.TypeString,
			ConflictsWith: []string{"stage"},
		},

		"stage": {
			Optional:      true,
			Type:          pluginsdk.TypeList,
			ConflictsWith: []string{"fleet_update_strategy_id"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Required:     true,
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"group": {
						Required: true,
						Type:     pluginsdk.TypeList,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Required:     true,
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"after_stage_wait_in_seconds": {
						Optional: true,
						Type:     pluginsdk.TypeInt,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
