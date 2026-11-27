// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_hyperv_network_mapping

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (s Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description:  "The name of the HyperV network mapping.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"recovery_vault_id": {
			Description:  "The ID of the Recovery Services Vault where the HyperV network mapping should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: vaults.ValidateVaultID,
		},

		"source_system_center_virtual_machine_manager_name": {
			Description:  "The name of source System Center Virtual Machine Manager where the source network exists.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"source_network_name": {
			Description:      "The Name of the primary network.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     validation.StringIsNotEmpty,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"target_network_id": {
			Description:      "The ID of the recovery network.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     azure.ValidateResourceID,
			DiffSuppressFunc: suppress.CaseDifference,
		},
	}
}

func (s Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}
