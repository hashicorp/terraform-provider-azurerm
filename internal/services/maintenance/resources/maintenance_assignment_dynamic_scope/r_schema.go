// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_assignment_dynamic_scope

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/configurationassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"maintenance_configuration_id": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     maintenanceconfigurations.ValidateMaintenanceConfigurationID,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"filter": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"locations": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						AtLeastOneOf: []string{"filter.0.locations", "filter.0.os_types", "filter.0.resource_groups", "filter.0.resource_types", "filter.0.tags"},
					},
					"os_types": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"Linux",
								"Windows",
							}, true),
						},
						AtLeastOneOf: []string{"filter.0.locations", "filter.0.os_types", "filter.0.resource_groups", "filter.0.resource_types", "filter.0.tags"},
					},
					"resource_groups": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						AtLeastOneOf: []string{"filter.0.locations", "filter.0.os_types", "filter.0.resource_groups", "filter.0.resource_types", "filter.0.tags"},
					},
					"resource_types": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"Microsoft.Compute/virtualMachines",
								"Microsoft.HybridCompute/machines",
							}, true),
						},
						AtLeastOneOf: []string{"filter.0.locations", "filter.0.os_types", "filter.0.resource_groups", "filter.0.resource_types", "filter.0.tags"},
					},
					"tags": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"tag": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"values": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
						AtLeastOneOf: []string{"filter.0.locations", "filter.0.os_types", "filter.0.resource_groups", "filter.0.resource_types", "filter.0.tags"},
					},

					"tag_filter": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      configurationassignments.TagOperatorsAny,
						ValidateFunc: validation.StringInSlice(configurationassignments.PossibleValuesForTagOperators(), true),
						RequiredWith: []string{
							"filter.0.tags",
						},
					},
				},
			},
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
