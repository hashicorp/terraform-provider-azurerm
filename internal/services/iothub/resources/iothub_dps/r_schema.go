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
			Description:  "The name of the Iot Device Provisioning Service resource.",
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
						Description: "The name of the SKU. Currently can only be set to `S1`.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(devices.IotHubSkuS1),
						}, false),
					},

					"capacity": {
						Description:  "The number of provisioned IoT Device Provisioning Service units.",
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
						Description:  "The connection string to connect to the IoT Hub.",
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
						Description:  "The location of the IoT hub.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						StateFunc:    location.StateFunc,
					},
					"apply_allocation_policy": {
						Description: "Whether to apply allocation policies to the IoT Hub.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"allocation_weight": {
						Description:  "The weight applied to the IoT Hub.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntBetween(0, 1000),
					},
					"hostname": {
						Description: "The IoT Hub hostname.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
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
						Description:  "The name of the filter.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"ip_mask": {
						Description:  "The IP address range in CIDR notation for the rule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsCIDRIPv4,
					},
					"action": {
						Description:  "The desired action for requests captured by this rule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInEnumSlice(devices.PossibleIPFilterActionTypeValues(), false),
					},
					"target": {
						Description:  "Target for requests captured by this rule.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(iotdpsresource.PossibleValuesForIPFilterTargetType(), false),
					},
				},
			},
		},

		"allocation_policy": {
			Description:  "The allocation policy of the IoT Device Provisioning Service (`Hashed`, `GeoLatency` or `Static`).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(iotdpsresource.AllocationPolicyHashed),
			ValidateFunc: validation.StringInSlice(iotdpsresource.PossibleValuesForAllocationPolicy(), false),
		},

		"data_residency_enabled": {
			Description: "Whether the IoT Device Provisioning Service has data residency enabled, removing the cross geo-pair disaster recovery.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"public_network_access_enabled": {
			Description: "Whether requests from Public Network are allowed.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"device_provisioning_host_name": {
			Description: "The device endpoint of the IoT Device Provisioning Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"id_scope": {
			Description: "The unique identifier of the IoT Device Provisioning Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"service_operations_host_name": {
			Description: "The service endpoint of the IoT Device Provisioning Service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
