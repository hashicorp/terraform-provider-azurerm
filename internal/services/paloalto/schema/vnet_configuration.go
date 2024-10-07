// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VnetConfiguration struct {
	VNetID            string `tfschema:"virtual_network_id"`
	TrustedSubnetID   string `tfschema:"trusted_subnet_id"`
	UntrustedSubnetID string `tfschema:"untrusted_subnet_id"`
	IpOfTrust         string `tfschema:"ip_of_trust_for_user_defined_routes"` // TODO - What is this?
}

func VnetConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"virtual_network_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: commonids.ValidateVirtualNetworkID,
				},

				"trusted_subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: commonids.ValidateSubnetID,
				},

				"untrusted_subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: commonids.ValidateSubnetID,
				},

				"ip_of_trust_for_user_defined_routes": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}
