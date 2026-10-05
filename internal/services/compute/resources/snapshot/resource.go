// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package snapshot

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/snapshots"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/resources/snapshot/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSnapshotCreateUpdate,
		Read:   resourceSnapshotRead,
		Update: resourceSnapshotCreateUpdate,
		Delete: resourceSnapshotDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := snapshots.ParseSnapshotID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SnapshotV0ToV1{},
		}),

		Schema: snapshotSchema(),

		// Encryption Settings cannot be disabled once enabled
		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("encryption_settings", func(ctx context.Context, old, new, meta any) bool {
				return len(old.([]any)) > 0 && len(new.([]any)) == 0
			}),
		),
	}
}
