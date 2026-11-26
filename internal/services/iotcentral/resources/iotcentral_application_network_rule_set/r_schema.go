// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iotcentral_application_network_rule_set

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"iotcentral_application_id": {
			Description:  "The ID of the IoT Central Application.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apps.ValidateIotAppID,
		},

		"apply_to_device": {
			Description: "Whether these IP Rules apply for device connectivity to IoT Hub and Device Provisioning Service associated with this IoT Central Application.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"default_action": {
			Description:  "The default action for the IoT Central Application Network Rule Set.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(apps.NetworkActionDeny),
			ValidateFunc: validation.StringInSlice(apps.PossibleValuesForNetworkAction(), false),
		},

		"ip_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the IP Rule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.IoTHubIpRuleName,
					},
					"ip_mask": {
						Description:  "The IP address range in CIDR notation for the IP Rule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsCIDRIPv4,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
