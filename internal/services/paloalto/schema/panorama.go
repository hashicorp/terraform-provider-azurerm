package schema

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Panorama struct {
	Name string `tfschema:"name,omitempty"`

	// Computed
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

				// Computed - Parsed out from the b64 config string
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true, // TODO - Check this contained in the b64 config?
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
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}
