// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job_schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := streamanalyticsClient.From(metadata.Client).JobsClient
			id, err := parse.StreamingJobScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			streamingJobId := streamingjobs.NewStreamingJobID(id.SubscriptionId, id.ResourceGroup, id.StreamingJobName)
			if err := client.StopThenPoll(ctx, streamingJobId); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}
