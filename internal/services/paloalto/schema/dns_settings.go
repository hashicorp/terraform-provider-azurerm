package schema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DNSSettings struct {
	DnsServers []string `tfschema:"dns_servers"`
	AzureDNS   bool     `tfschema:"use_azure_dns"`
}

func DNSSettingsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dns_servers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 2,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validate.IPv4Address,
					},
					ConflictsWith: []string{
						"dns_settings.0.use_azure_dns",
					},
				},

				"use_azure_dns": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
					ConflictsWith: []string{
						"dns_settings.0.dns_servers",
					},
				},
			},
		},
	}
}
