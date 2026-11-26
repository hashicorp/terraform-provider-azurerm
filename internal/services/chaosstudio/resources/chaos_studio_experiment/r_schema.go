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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"selectors": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"chaos_studio_target_ids": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: commonids.ValidateChaosStudioTargetID,
						},
					},
				},
			},
		},
		"steps": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"branch": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"actions": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"action_type": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													continuousActionType,
													delayActionType,
													discreteActionType,
												}, false),
											},
											// the different types of actions require different properties to be set
											// which is why the validation for these is done in expandActions
											"urn": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
											"selector_name": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
											"duration": {
												Type:     pluginsdk.TypeString,
												Optional: true,
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
