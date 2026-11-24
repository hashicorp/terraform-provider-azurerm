// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_protected_vm

import (
	"slices"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceRecoveryServicesBackupProtectedVMSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"source_vm_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
			ValidateFunc: validation.Any(
				validation.StringIsEmpty,
				azure.ValidateResourceID,
			),
			// TODO: make this case sensitive once the API's fixed https://github.com/Azure/azure-rest-api-specs/issues/10357
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"backup_policy_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: protectionpolicies.ValidateBackupPolicyID,
		},

		"exclude_disk_luns": {
			Type:          pluginsdk.TypeSet,
			ConflictsWith: []string{"include_disk_luns"},
			Optional:      true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeInt,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},

		"include_disk_luns": {
			Type:          pluginsdk.TypeSet,
			ConflictsWith: []string{"exclude_disk_luns"},
			Optional:      true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeInt,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},

		"protection_state": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because `protection_state` is set by Azure and may not be a persistent value.
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				// While not a persistent state, `Protected` is an option to allow a path from `BackupsSuspended`/`ProtectionStopped` to a protected state.
				string(protecteditems.ProtectionStateProtected),
				string(protecteditems.ProtectionStateBackupsSuspended),
				string(protecteditems.ProtectionStateProtectionStopped),
			}, false),
			DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
				// We suppress the diff if the only change is from "IRPending" or "ProtectionPaused" to "Protected".
				// These states are not persistent and are set by Azure based on the current protection state.
				// While `Invalid` and `ProtectionError` are also not configurable, we're opting to output this in the diff
				// as these states should indicate to the user that there is an error with the backup protected VM resource requiring attention.
				suppressStates := []string{
					string(protecteditems.ProtectedItemStateIRPending),
					string(protecteditems.ProtectedItemStateProtectionPaused),
				}

				if new == string(protecteditems.ProtectionStateProtected) && slices.Contains(suppressStates, old) {
					return true
				}

				return false
			},
		},
	}
}
