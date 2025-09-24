// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"slices"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2025-02-01/protecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2025-02-01/protectionpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func BackupPolicyIdSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Optional:     true,
		ValidateFunc: protectionpolicies.ValidateBackupPolicyID,
	}
}

func RecoveryVaultNameSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validate.RecoveryServicesVaultName,
	}
}

func SourceVMIdSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Optional: true,
		Computed: true,
		ForceNew: true,
		ValidateFunc: validation.Any(
			validation.StringIsEmpty,
			azure.ValidateResourceID,
		),
		// TODO: make this case sensitive once the API's fixed https://github.com/Azure/azure-rest-api-specs/issues/10357
		DiffSuppressFunc: suppress.CaseDifference,
	}
}

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
