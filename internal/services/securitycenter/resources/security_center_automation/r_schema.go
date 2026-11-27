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
			Description:  "The name which should be used for this Security Center Automation.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": {
			Description: "The Azure Region where the Security Center Automation should exist.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			StateFunc:   location.StateFunc,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"enabled": {
			Description: "Boolean to enable or disable this Security Center Automation.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"description": {
			Description:  "The description for the Security Center Automation.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"scopes": {
			Description: "A list of scopes on which the automation logic is applied, at least one is required. Supported scopes are a subscription (in this format `/subscriptions/00000000-0000-0000-0000-000000000000`) or a Resource Group under that subscription (in the format `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example`). The automation will only apply on defined scopes.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: commonids.ValidateScopeID,
			},
		},

		"action": {
			Description: "An `action` tells this automation where the data is to be sent to upon being evaluated by the rules in the `source`.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MinItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Description:  "Type of Azure resource to send data to.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(automations.PossibleValuesForActionType(), false),
					},

					"resource_id": {
						Description:  "The resource ID of the target Logic App, Event Hub namespace or Log Analytics Workspace.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"trigger_url": {
						Description:  "The callback URL to trigger the Logic App that will receive and process data sent by this automation. This can be found in the Azure Portal under \"See trigger history\".",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"connection_string": {
						Description:  "A connection string to send data to the target Event Hub namespace, this should include a key with send permissions.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"source": {
			Description: "A `source` defines what data types will be processed and a set of rules to filter that data.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MinItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"event_source": {
						Description:  "Type of data that will trigger this automation. Note: Assessments are also referred to as recommendations.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(automations.PossibleValuesForEventSource(), false),
					},

					"rule_set": {
						Description: "A set of rules which evaluate upon event and data interception.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"rule": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"property_path": {
												Description: "The JPath of the entity model property that should be checked.",
												Type:        pluginsdk.TypeString,
												Required:    true,
											},
											"expected_value": {
												Description: "A value that will be compared with the value in `property_path`.",
												Type:        pluginsdk.TypeString,
												Required:    true,
											},
											"operator": {
												Description:  "The comparison operator to use, must be one of: `Contains`, `EndsWith`, `Equals`, `GreaterThan`, `GreaterThanOrEqualTo`, `LesserThan`, `LesserThanOrEqualTo`, `NotEquals`, `StartsWith`.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(automations.PossibleValuesForOperator(), false),
											},
											"property_type": {
												Description:  "The data type of the compared operands, must be one of: `Integer`, `String`, `Boolean` or `Number`.",
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
