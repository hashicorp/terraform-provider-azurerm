// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_cloud_serverless

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Elastic.ServerlessMonitorClient
			id, err := elasticmonitorresources.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.MonitorsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return nil
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if err := validateElasticCloudServerlessMonitor(id, existing.Model); err != nil {
				return err
			}

			if err := client.MonitorsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
