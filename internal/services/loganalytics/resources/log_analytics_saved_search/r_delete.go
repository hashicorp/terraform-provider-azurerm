// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_saved_search

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/savedsearches"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	loganalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsSavedSearchDelete(d *pluginsdk.ResourceData, meta any) error {
	client := loganalyticsClient.From(meta).SavedSearchesClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()
	id, err := savedsearches.ParseSavedSearchID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
