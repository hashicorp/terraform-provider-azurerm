// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_linux_virtual_machine

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func devTestLinuxVirtualMachineSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Dev Test Machine.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevTestVirtualMachineName(62),
		},

		"lab_name": {
			Description:  "The name of the Dev Test Lab in which the Virtual Machine should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevTestLabName(),
		},

		// There's a bug in the Azure API where this is returned in lower-case
		// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"location": commonschema.Location(),

		"size": {
			Description: "The Machine Size to use for this Virtual Machine, such as `Standard_D4_v5`.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			// since this isn't returned from the API
			ForceNew: true,
		},

		"username": {
			Description: "The Username associated with the local administrator on this Virtual Machine.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"storage_type": {
			Description: "The type of Storage to use on this Virtual Machine.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"Standard",
				"Premium",
			}, false),
		},

		"lab_subnet_name": {
			Description: "The name of a Subnet within the Dev Test Virtual Network where this machine should exist.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			// since this isn't returned from the API
			ForceNew: true,
		},

		"lab_virtual_network_id": {
			Description: "The ID of the Dev Test Virtual Network where this Virtual Machine should be created.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			// since this isn't returned from the API
			ForceNew: true,
		},

		"allow_claim": {
			Description: "Can this Virtual Machine be claimed by users?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"disallow_public_ip_address": {
			Description: "Should the Virtual Machine be created without a Public IP Address?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"password": {
			Description: "The Password associated with the `username` used to login to this Virtual Machine.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// since this isn't returned from the API
			ForceNew:  true,
			Sensitive: true,
		},

		"ssh_key": {
			Description: "The SSH Key associated with the `username` used to login to this Virtual Machine.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// since this isn't returned from the API
			ForceNew: true,
		},

		"gallery_image_reference": helpers.SchemaDevTestVirtualMachineGalleryImageReference(),

		"inbound_nat_rule": helpers.SchemaDevTestVirtualMachineInboundNatRule(),

		"notes": {
			Description: "Any notes about the Virtual Machine.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"tags": commonschema.Tags(),

		"fqdn": {
			Description: "The FQDN of the Virtual Machine.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"unique_identifier": {
			Description: "The unique immutable identifier of the Virtual Machine.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
