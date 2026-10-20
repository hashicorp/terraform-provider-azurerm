// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_fleet_update_strategy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01/fleetupdatestrategies"
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

		"kubernetes_fleet_manager_id": commonschema.ResourceIDReferenceRequiredForceNew(&fleetupdatestrategies.FleetId{}),

		"stage": {
			Required: true,
			Type:     pluginsdk.TypeList,
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
