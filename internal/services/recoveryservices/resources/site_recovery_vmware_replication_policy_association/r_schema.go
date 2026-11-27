// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replication_policy_association

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (s Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the replication policy association.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"recovery_vault_id": {
			Description:  "The ID of the Recovery Service Vault to which the policy should be associated.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: vaults.ValidateVaultID,
		},

		"policy_id": {
			Description:  "The ID of the VMWare replication policy which to be associated.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationpolicies.ValidateReplicationPolicyID,
		},
	}
}

func (s Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
