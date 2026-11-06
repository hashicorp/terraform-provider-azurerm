// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_automation

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2019-01-01-preview/automations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func securityCenterAutomationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			ForceNew:  true,
			StateFunc: location.StateFunc,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Default:  true,
			Optional: true,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"scopes": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateScopeID,
			},
		},

		"action": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(automations.PossibleValuesForActionType(), false),
					},

					"resource_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"trigger_url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"connection_string": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"source": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"event_source": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(automations.PossibleValuesForEventSource(), false),
					},

					"rule_set": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"rule": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"property_path": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"expected_value": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"operator": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(automations.PossibleValuesForOperator(), false),
											},
											"property_type": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(automations.PossibleValuesForPropertyType(), false),
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

		"tags": commonschema.Tags(),
	}
}
