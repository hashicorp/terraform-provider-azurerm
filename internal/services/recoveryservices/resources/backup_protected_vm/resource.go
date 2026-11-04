// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_protected_vm

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRecoveryServicesBackupProtectedVMCreate,
		Read:   resourceRecoveryServicesBackupProtectedVMRead,
		Update: resourceRecoveryServicesBackupProtectedVMUpdate,
		Delete: resourceRecoveryServicesBackupProtectedVMDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := protecteditems.ParseProtectedItemID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(120 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(120 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(80 * time.Minute),
		},

		Schema: resourceRecoveryServicesBackupProtectedVMSchema(),

		// It's possible to remove the associated vm from the protected backup so we'll only ForceNew this attribute if it's
		// changing to something other than empty.
		CustomizeDiff: pluginsdk.ForceNewIfChange("source_vm_id", func(ctx context.Context, old, new, meta any) bool {
			return new.(string) != "" && old.(string) != new.(string)
		}),
	}
}
