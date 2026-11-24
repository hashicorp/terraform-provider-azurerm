// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package hdinsight_hadoop_cluster

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SchemaHDInsightsHttpsEndpoints() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description: "The HTTPS Connectivity Endpoint for this HDInsight Hadoop Cluster.",
		Type:        pluginsdk.TypeList,
		Optional:    true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"access_modes": {
					Description: "A list of access modes for the application.",
					Type:        pluginsdk.TypeList,
					Optional:    true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"destination_port": {
					Description:  "The destination port to connect to.",
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IsPortNumber,
				},

				"disable_gateway_auth": {
					Description: "The value indicates whether the gateway authentication is enabled or not.",
					Type:        pluginsdk.TypeBool,
					Optional:    true,
				},

				"private_ip_address": {
					Description:  "The private IP address of the endpoint.",
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsIPAddress,
				},

				"sub_domain_suffix": {
					Description:  "The application's subdomain suffix.",
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}
