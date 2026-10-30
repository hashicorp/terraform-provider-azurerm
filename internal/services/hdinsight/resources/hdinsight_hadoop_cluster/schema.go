// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package hdinsight_hadoop_cluster

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SchemaHDInsightsHttpsEndpoints() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"access_modes": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"destination_port": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IsPortNumber,
				},

				"disable_gateway_auth": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"private_ip_address": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},

				"sub_domain_suffix": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}
