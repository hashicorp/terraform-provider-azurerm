// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_virtual_machine_configuration_assignment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/guestconfiguration/2024-04-05/guestconfigurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourcePolicyVirtualMachineConfigurationAssignmentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Guest Configuration that will be assigned in this Guest Configuration Assignment.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"virtual_machine_id": {
			Description:  "The resource ID of the Policy Virtual Machine which this Guest Configuration Assignment should apply to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"assignment_type": {
						Description:  "The assignment type for the Guest Configuration Assignment.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(guestconfigurationassignments.PossibleValuesForAssignmentType(), false),
					},

					"content_hash": {
						Description:  "The content hash for the Guest Configuration package.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"content_uri": {
						Description:  "The content URI where the Guest Configuration package is stored.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
					},

					"parameter": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description:  "The name of the configuration parameter to check.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"value": {
									Description: "The value to check the configuration parameter with.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
							},
						},
					},

					"version": {
						Description: "The version of the Guest Configuration that will be assigned in this Guest Configuration Assignment.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
				},
			},
		},
	}
}
