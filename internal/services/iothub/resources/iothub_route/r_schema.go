// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_route

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func iothubRouteSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the route.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
				"Route Name name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"iothub_name": {
			Description:  "The name of the IoTHub to which this Route belongs.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IoTHubName,
		},

		"source": {
			Description: "The source that the routing rule is to be applied to.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(devices.RoutingSourceDeviceConnectionStateEvents),
				string(devices.RoutingSourceDeviceJobLifecycleEvents),
				string(devices.RoutingSourceDeviceLifecycleEvents),
				string(devices.RoutingSourceDeviceMessages),
				string(devices.RoutingSourceDigitalTwinChangeEvents),
				string(devices.RoutingSourceInvalid),
				string(devices.RoutingSourceTwinChangeEvents),
			}, false),
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
			Description: "The list of endpoints to which messages that satisfy the condition are routed. Currently only one endpoint is allowed.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			// Currently only one endpoint is allowed. With that comment from Microsoft, we'll leave this open to enhancement when they add multiple endpoint support.
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"enabled": {
			Description: "Whether a route is enabled.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},
	}
}
