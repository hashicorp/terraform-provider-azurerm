// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_windows_virtual_machine

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func devTestWindowsVirtualMachineSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevTestVirtualMachineName(15),
		},

		"lab_name": {
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
			Type:     pluginsdk.TypeString,
			Required: true,
			// since this isn't returned from the API
			ForceNew: true,
		},

		"username": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"password": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// since this isn't returned from the API
			ForceNew:  true,
			Sensitive: true,
		},

		"storage_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Standard",
				"Premium",
			}, false),
		},

		"lab_subnet_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// since this isn't returned from the API
			ForceNew: true,
		},

		"lab_virtual_network_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// since this isn't returned from the API
			ForceNew: true,
		},

		"allow_claim": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"disallow_public_ip_address": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"gallery_image_reference": helpers.SchemaDevTestVirtualMachineGalleryImageReference(),

		"inbound_nat_rule": helpers.SchemaDevTestVirtualMachineInboundNatRule(),

		"notes": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": commonschema.Tags(),

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"unique_identifier": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
