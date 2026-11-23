// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/resources/stream_analytics_job/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsJobCreate,
		Read:   resourceStreamAnalyticsJobRead,
		Update: resourceStreamAnalyticsJobUpdate,
		Delete: resourceStreamAnalyticsJobDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := streamingjobs.ParseStreamingJobID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsJobV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: streamAnalyticsJobSchema(),

		CustomizeDiff: func(ctx context.Context, d *pluginsdk.ResourceDiff, i any) error {
			if d.Get("job_storage_account.0.authentication_mode") == string(streamingjobs.AuthenticationModeMsi) && d.Get("job_storage_account.0.account_key") != "" {
				return fmt.Errorf("`job_storage_account.0.account_key` cannot be set when `job_storage_account.0.authentication_mode` is `Msi`")
			}
			return nil
		},
	}
}
