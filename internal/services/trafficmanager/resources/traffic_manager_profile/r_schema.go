// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_profile

import (
	"math"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func trafficManagerProfileSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"traffic_routing_method": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(profiles.PossibleValuesForTrafficRoutingMethod(), false),
		},

		"dns_config": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"relative_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
					"ttl": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, math.MaxInt32),
					},
				},
			},
		},

		"monitor_config": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expected_status_code_ranges": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.StatusCodeRange,
						},
					},

					"custom_header": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"value": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},

					"protocol": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(profiles.PossibleValuesForMonitorProtocol(), false),
					},

					"port": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 65535),
					},

					"path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"interval_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      30,
						ValidateFunc: validation.IntInSlice([]int{10, 30}),
					},

					"timeout_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      10,
						ValidateFunc: validation.IntBetween(5, 10),
					},

					"tolerated_number_of_failures": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      3,
						ValidateFunc: validation.IntBetween(0, 9),
					},
				},
			},
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"max_return": { // azignore:AZP003 - named `maximum_return` in the data source to follow new naming conventions
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 8),
		},

		"traffic_view_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}

	if !features.SixPointOh() {
		s["profile_status"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ValidateFunc:  validation.StringInSlice(profiles.PossibleValuesForProfileStatus(), false),
			ConflictsWith: []string{"enabled"},
			Deprecated:    "`profile_status` has been deprecated in favour of `enabled` and will be removed in v6.0 of the AzureRM Provider",
		}

		s["enabled"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"profile_status"},
		}
	}

	return s
}
