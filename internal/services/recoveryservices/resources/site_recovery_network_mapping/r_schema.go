// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_network_mapping

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func siteRecoveryNetworkMappingSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the network mapping.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupName(),

		"recovery_vault_name": {
			Description:  "The name of the vault that should be updated.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},
		"source_recovery_fabric_name": {
			Description:  "The ASR fabric where mapping should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"target_recovery_fabric_name": {
			Description:  "The Azure Site Recovery fabric object corresponding to the recovery Azure region.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"source_network_id": {
			Description:      "The ID of the primary network.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     azure.ValidateResourceID,
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
