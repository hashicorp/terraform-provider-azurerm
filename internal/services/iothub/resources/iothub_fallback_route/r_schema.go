// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_fallback_route

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func iothubFallbackRouteSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"iothub_name": {
			Description:  "The name of the IoTHub to which this Fallback Route belongs.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IoTHubName,
		},

		"source": {
			Description: "The source that the routing rule is to be applied to.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(devices.RoutingSourceDeviceConnectionStateEvents),
				string(devices.RoutingSourceDeviceJobLifecycleEvents),
				string(devices.RoutingSourceDeviceLifecycleEvents),
				string(devices.RoutingSourceDeviceMessages),
				string(devices.RoutingSourceDigitalTwinChangeEvents),
				string(devices.RoutingSourceInvalid),
				string(devices.RoutingSourceTwinChangeEvents),
			}, false),
			Default: devices.RoutingSourceDeviceMessages,
		},

		"condition": {
			Description: "The condition that is evaluated to apply the routing rule. For grammar, see: <https://learn.microsoft.com/azure/iot-hub/iot-hub-devguide-query-language>.",
			// The condition is a string value representing device-to-cloud message routes query expression
			// https://learn.microsoft.com/azure/iot-hub/iot-hub-devguide-query-language#device-to-cloud-message-routes-query-expressions
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "true",
		},

		"endpoint_names": {
			Description: "The endpoints to which messages that satisfy the condition are routed. Currently only 1 endpoint is allowed.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			// Currently only one endpoint is allowed. With that comment from Microsoft, we'll leave this open to enhancement when they add multiple endpoint support.
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.IoTHubEndpointName,
			},
		},

		"enabled": {
			Description: "Used to specify whether the fallback route is enabled.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},
	}
}
