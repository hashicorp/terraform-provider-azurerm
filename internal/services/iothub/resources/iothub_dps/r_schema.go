// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func iothubDpsSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IoTHubName,
		},

		"resource_group_name": commonschema.ResourceGroupName(), // azure.SchemaResourceGroupNameDiffSuppress(),

		"location": commonschema.Location(),

		"sku": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(devices.IotHubSkuS1),
						}, false),
					},

					"capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 200),
					},
				},
			},
		},

		"linked_hub": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"connection_string": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						// Azure returns the key as ****. We'll suppress that here.
						DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
							secretKeyRegex := regexp.MustCompile("(SharedAccessKey)=[^;]+")
							maskedNew := secretKeyRegex.ReplaceAllString(new, "$1=****")
							return (new == d.Get(k).(string)) && (maskedNew == old)
						},
						Sensitive: true,
					},
					"location": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						StateFunc:    location.StateFunc,
					},
					"apply_allocation_policy": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"allocation_weight": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
					"hostname": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"ip_filter_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"ip_mask": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsCIDRIPv4,
					},
					"action": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInEnumSlice(devices.PossibleIPFilterActionTypeValues(), false),
					},
					"target": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(iotdpsresource.PossibleValuesForIPFilterTargetType(), false),
					},
				},
			},
		},

		"allocation_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(iotdpsresource.AllocationPolicyHashed),
			ValidateFunc: validation.StringInSlice(iotdpsresource.PossibleValuesForAllocationPolicy(), false),
		},

		"data_residency_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"device_provisioning_host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"id_scope": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"service_operations_host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
