// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replication_policy

import (
	"context"
	"errors"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func siteRecoveryReplicationPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the replication policy.",
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
		"recovery_point_retention_in_minutes": {
			Description:  "The duration in minutes for which the recovery points need to be stored.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(0, 365*24*60),
		},
		"application_consistent_snapshot_frequency_in_minutes": {
			Description:  "The frequency(in minutes) at which to create application consistent recovery points.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(0, 365*24*60),
		},
	}
}

func resourceSiteRecoveryReplicationPolicyCustomDiff(ctx context.Context, d *pluginsdk.ResourceDiff, i any) error {
	retention := d.Get("recovery_point_retention_in_minutes").(int)
	frequency := d.Get("application_consistent_snapshot_frequency_in_minutes").(int)

	if retention == 0 && frequency > 0 {
		return errors.New("application_consistent_snapshot_frequency_in_minutes cannot be greater than zero when recovery_point_retention_in_minutes is set to zero")
	}

	return nil
}
