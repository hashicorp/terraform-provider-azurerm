// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshotpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetAppSnapshotPolicyCreate,
		Read:   resourceNetAppSnapshotPolicyRead,
		Update: resourceNetAppSnapshotPolicyUpdate,
		Delete: resourceNetAppSnapshotPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := snapshotpolicies.ParseSnapshotPolicyID(id)
			return err
		}),

		Schema: netappSnapshotPolicySchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("hourly_schedule", func(ctx context.Context, old, new, meta any) bool {
				return len(old.([]any)) > 0 && len(new.([]any)) == 0
			}),

			pluginsdk.ForceNewIfChange("daily_schedule", func(ctx context.Context, old, new, meta any) bool {
				return len(old.([]any)) > 0 && len(new.([]any)) == 0
			}),

			pluginsdk.ForceNewIfChange("weekly_schedule", func(ctx context.Context, old, new, meta any) bool {
				return len(old.([]any)) > 0 && len(new.([]any)) == 0
			}),

			pluginsdk.ForceNewIfChange("monthly_schedule", func(ctx context.Context, old, new, meta any) bool {
				return len(old.([]any)) > 0 && len(new.([]any)) == 0
			}),
		),
	}
}
