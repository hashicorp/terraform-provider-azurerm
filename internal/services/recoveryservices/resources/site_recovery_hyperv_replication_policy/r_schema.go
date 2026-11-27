// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_hyperv_replication_policy

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description:  "The name of the replication policy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"recovery_vault_id": {
			Description:  "The ID of the vault that should be updated.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: vaults.ValidateVaultID,
		},
		"recovery_point_retention_in_hours": {
			Description:  "The duration in hours for which the recovery points need to be stored.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(0, 24),
		},
		"application_consistent_snapshot_frequency_in_hours": {
			Description:  "The frequency at which to create application consistent recovery points.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(0, 12),
		},
		"replication_interval_in_seconds": {
			Description:  "How frequently data should be synchronized between source and target locations.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntInSlice([]int{30, 300}),
		},
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}
