// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/dataexport"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &LogAnalyticsDataExportPoller{}

type LogAnalyticsDataExportPoller struct {
	client *dataexport.DataExportClient
	id     dataexport.DataExportId
}

var (
	pollingSuccess = pollers.PollResult{
		Status:       pollers.PollingStatusSucceeded,
		PollInterval: 10 * time.Second,
	}
	pollingInProgress = pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}
)

func NewLogAnalyticsDataExportPoller(client *dataexport.DataExportClient, id dataexport.DataExportId) *LogAnalyticsDataExportPoller {
	return &LogAnalyticsDataExportPoller{
		client: client,
		id:     id,
	}
}

func (p LogAnalyticsDataExportPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollingInProgress, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}
	return &pollingSuccess, nil
}
