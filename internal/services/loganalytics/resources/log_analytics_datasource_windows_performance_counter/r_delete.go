// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_datasource_windows_performance_counter

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	loganalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsDataSourceWindowsPerformanceCounterDelete(d *pluginsdk.ResourceData, meta any) error {
	client := loganalyticsClient.From(meta).DataSourcesClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := datasources.ParseDataSourceID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("failed to delete Windows Performance Counter %s: %+v", *id, err)
	}

	return nil
}
