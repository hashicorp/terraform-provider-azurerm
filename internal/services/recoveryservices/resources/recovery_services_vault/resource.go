// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRecoveryServicesVaultCreate,
		Read:   resourceRecoveryServicesVaultRead,
		Update: resourceRecoveryServicesVaultUpdate,
		Delete: resourceRecoveryServicesVaultDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := vaults.ParseVaultID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(120 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: recoveryServicesVaultSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("cross_region_restore_enabled", func(ctx context.Context, old, new, meta any) bool {
				return old.(bool) && !new.(bool)
			}),
			pluginsdk.ForceNewIfChange("immutability", func(ctx context.Context, old, new, meta any) bool {
				return old.(string) == string(vaults.ImmutabilityStateLocked)
			}),
		),
	}
}
