// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_virtual_network

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func devTestVirtualNetworkSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Dev Test Virtual Network.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: ValidateDevTestVirtualNetworkName(),
		},

		"lab_name": {
			Description:  "The name of the Dev Test Lab in which the Virtual Network should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevTestLabName(),
		},

		// There's a bug in the Azure API where this is returned in lower-case
		// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"description": {
			Description: "A description for the Virtual Network.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"subnet": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// NOTE: O+C We generate the name and resource ID for this block so this should remain
			Computed: true,
			// whilst the API accepts multiple, in practice only one is usable
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name of the Subnet for this Virtual Network.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"use_in_virtual_machine_creation": {
						Description:  "Whether this subnet can be used for creating Virtual Machines.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(virtualnetworks.UsagePermissionTypeAllow),
						ValidateFunc: validate.DevTestVirtualNetworkUsagePermissionType(),
					},

					"use_public_ip_address": {
						Description:  "Whether Virtual Machines in this Subnet can use Public IP Addresses.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(virtualnetworks.UsagePermissionTypeAllow),
						ValidateFunc: validate.DevTestVirtualNetworkUsagePermissionType(),
					},

					"shared_public_ip_address": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allowed_ports": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"backend_port": {
												Description: "The port on the Virtual Machine that the traffic will be sent to.",
												Type:        pluginsdk.TypeInt,
												Optional:    true,
											},

											"transport_protocol": {
												Description:  "The transport protocol that the traffic will use.",
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringInSlice(virtualnetworks.PossibleValuesForTransportProtocol(), false),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"unique_identifier": {
			Description: "The unique immutable identifier of the Dev Test Virtual Network.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
