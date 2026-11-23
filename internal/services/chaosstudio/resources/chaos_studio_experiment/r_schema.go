// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package chaos_studio_experiment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"selectors": {
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
					"chaos_studio_target_ids": {
						Required: true,
						Type:     pluginsdk.TypeList,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: commonids.ValidateChaosStudioTargetID,
						},
					},
				},
			},
		},
		"steps": {
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
					"branch": {
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
								"actions": {
									Required: true,
									Type:     pluginsdk.TypeList,
									MinItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"action_type": {
												Required: true,
												Type:     pluginsdk.TypeString,
												ValidateFunc: validation.StringInSlice([]string{
													continuousActionType,
													delayActionType,
													discreteActionType,
												}, false),
											},
											// the different types of actions require different properties to be set
											// which is why the validation for these is done in expandActions
											"urn": {
												Optional: true,
												Type:     pluginsdk.TypeString,
											},
											"selector_name": {
												Optional: true,
												Type:     pluginsdk.TypeString,
											},
											"duration": {
												Optional: true,
												Type:     pluginsdk.TypeString,
											},
											"parameters": {
												Type:     pluginsdk.TypeMap,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
