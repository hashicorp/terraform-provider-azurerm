// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Panorama struct {
	Name            string `tfschema:"name"`
	DeviceGroupName string `tfschema:"device_group_name"`
	HostName        string `tfschema:"host_name"`
	PanoramaServer  string `tfschema:"panorama_server_1"`
	PanoramaServer2 string `tfschema:"panorama_server_2"`
	TplName         string `tfschema:"template_name"`
	VMAuthKey       string `tfschema:"virtual_machine_ssh_key"`
}

func PanoramaSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Description: "The name of this Firewall Resource.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"device_group_name": {
					Description: "The Device Group Name to which this Firewall Resource is registered.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"host_name": {
					Description: "The Host Name of this Firewall Resource.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"panorama_server_1": {
					Description: "The name of the First Panorana server.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"panorama_server_2": {
					Description: "The name of the Second Panorana server.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"template_name": {
					Description: "The name of the Panorama Template applied to this Firewall Resource.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},

				"virtual_machine_ssh_key": {
					Description: "The SSH Key to connect to the Firewall Resource.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
				},
			},
		},
	}
}
