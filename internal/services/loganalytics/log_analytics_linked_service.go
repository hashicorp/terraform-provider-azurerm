// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func logAnalyticsLinkedServiceDeleteWaitForState(ctx context.Context, client *linkedservices.LinkedServicesClient, timeout time.Duration, id linkedservices.LinkedServiceId) *pluginsdk.StateChangeConf {
	return &pluginsdk.StateChangeConf{
		Pending:    []string{"Deleting"},
		Target:     []string{"Deleted"},
		MinTimeout: 30 * time.Second,
		Timeout:    timeout,
		Refresh:    logAnalyticsLinkedServiceRefresh(ctx, client, id),
	}
}

func logAnalyticsLinkedServiceRefresh(ctx context.Context, client *linkedservices.LinkedServicesClient, id linkedservices.LinkedServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[INFO] checking on state of %s", id)

		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return nil, "nil", fmt.Errorf("polling for the status of %s: %+v", id, err)
			}

			return resp, "Deleted", nil
		}

		// (@WodansSon) - The service returns status code 200 even if the resource does not exist
		// instead it returns an empty slice...
		if resp.Model == nil {
			return resp, "Deleted", nil
		}

		return resp, "Deleting", nil
	}
}
