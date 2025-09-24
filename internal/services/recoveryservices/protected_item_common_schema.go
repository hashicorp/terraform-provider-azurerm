// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"slices"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2025-02-01/protecteditems"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// ProtectionStateSchema returns the common schema for protection_state field used across backup resources
func ProtectionStateSchema(validStates []string) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		// Note: O+C because `protection_state` is set by Azure and may not be a persistent value.
		Computed: true,
		ValidateFunc: validation.StringInSlice(validStates, false),
		DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
			// We suppress the diff if the only change is from "IRPending" or "ProtectionPaused" to "Protected".
			// These states are not persistent and are set by Azure based on the current protection state.
			// While `Invalid` and `ProtectionError` are also not configurable, we're opting to output this in the diff
			// as these states should indicate to the user that there is an error with the backup protected resource requiring attention.
			suppressStates := []string{
				string(protecteditems.ProtectedItemStateIRPending),
				string(protecteditems.ProtectedItemStateProtectionPaused),
			}

			if new == string(protecteditems.ProtectionStateProtected) && slices.Contains(suppressStates, old) {
				return true
			}

			return false
		},
	}
}

// BackupProtectedVMProtectionStateSchema returns the protection state schema for VM backup resources
func BackupProtectedVMProtectionStateSchema() *pluginsdk.Schema {
	return ProtectionStateSchema([]string{
		// While not a persistent state, `Protected` is an option to allow a path from `BackupsSuspended`/`ProtectionStopped` to a protected state.
		string(protecteditems.ProtectionStateProtected),
		string(protecteditems.ProtectionStateBackupsSuspended),
		string(protecteditems.ProtectionStateProtectionStopped),
	})
}

// BackupProtectedVMWorkloadProtectionStateSchema returns the protection state schema for VM workload backup resources
func BackupProtectedVMWorkloadProtectionStateSchema() *pluginsdk.Schema {
	return ProtectionStateSchema([]string{
		// While not a persistent state, `Protected` is an option to allow a path from `ProtectionStopped` to a protected state.
		string(protecteditems.ProtectionStateProtected),
		string(protecteditems.ProtectionStateProtectionStopped),
	})
}
