// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance_configuration

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/maintenanceconfigurations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func maintenanceConfigurationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"scope": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(maintenanceconfigurations.MaintenanceScopeExtension),
				string(maintenanceconfigurations.MaintenanceScopeHost),
				string(maintenanceconfigurations.MaintenanceScopeInGuestPatch),
				string(maintenanceconfigurations.MaintenanceScopeOSImage),
				string(maintenanceconfigurations.MaintenanceScopeSQLDB),
				string(maintenanceconfigurations.MaintenanceScopeSQLManagedInstance),
			}, false),
		},

		"visibility": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(maintenanceconfigurations.VisibilityCustom),
			ValidateFunc: validation.StringInSlice([]string{
				string(maintenanceconfigurations.VisibilityCustom),
				// Creating public configurations doesn't appear to be supported, API returns `Public Maintenance Configuration must set correct properties`
				// string(maintenance.VisibilityPublic),
			}, false),
		},

		"window": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"start_date_time": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"expiration_date_time": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"duration": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$"),
							"duration must match the format HH:mm",
						),
					},
					"time_zone": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.MaintenanceTimeZone(),
					},
					"recur_every": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"install_patches": {
			Type:     pluginsdk.TypeList,
			MinItems: 1,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"linux": {
						Type:     pluginsdk.TypeList,
						MinItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"classifications_to_include": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"Critical",
											"Security",
											"Other",
										}, false),
									},
								},
								"package_names_mask_to_exclude": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"package_names_mask_to_include": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"windows": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"classifications_to_include": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"Critical",
											"Security",
											"UpdateRollup",
											"FeaturePack",
											"ServicePack",
											"Definition",
											"Tools",
											"Updates",
										}, false),
									},
								},
								"kb_numbers_to_exclude": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"kb_numbers_to_include": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"reboot": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(maintenanceconfigurations.PossibleValuesForRebootOptions(), false),
					},
				},
			},
		},

		"in_guest_user_patch_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Platform",
				"User",
			}, false),
		},

		"properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"tags": commonschema.Tags(),
	}
}
