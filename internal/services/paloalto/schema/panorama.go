// Copyright (c) HashiCorp, Inc.
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
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"device_group_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"host_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"panorama_server_1": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"panorama_server_2": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"template_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"virtual_machine_ssh_key": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},
			},
		},
	}
}
