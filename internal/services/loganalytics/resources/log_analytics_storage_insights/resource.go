// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_storage_insights

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/storageinsights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsStorageInsightsCreateUpdate,
		Read:   resourceLogAnalyticsStorageInsightsRead,
		Update: resourceLogAnalyticsStorageInsightsCreateUpdate,
		Delete: resourceLogAnalyticsStorageInsightsDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := storageinsights.ParseStorageInsightConfigID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta any) ([]*pluginsdk.ResourceData, error) {
			if v, ok := d.GetOk("storage_account_key"); ok && v.(string) != "" {
				d.Set("storage_account_key", v)
			}

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Schema: resourceLogAnalyticsStorageInsightsSchema(),
	}
}
