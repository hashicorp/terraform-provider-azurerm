package newrelicmonitorresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsResubscribeOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *NewRelicMonitorResource
}

// MonitorsResubscribe ...
func (c NewRelicMonitorResourcesClient) MonitorsResubscribe(ctx context.Context, id MonitorId, input ResubscribeProperties) (result MonitorsResubscribeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/resubscribe", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// MonitorsResubscribeThenPoll performs MonitorsResubscribe then polls until it's completed
func (c NewRelicMonitorResourcesClient) MonitorsResubscribeThenPoll(ctx context.Context, id MonitorId, input ResubscribeProperties) error {
	return c.MonitorsResubscribeCallbackThenPoll(ctx, id, input, nil)
}

// MonitorsResubscribeCallbackThenPoll performs MonitorsResubscribe, runs the optional callback function, then polls until it's completed
func (c NewRelicMonitorResourcesClient) MonitorsResubscribeCallbackThenPoll(ctx context.Context, id MonitorId, input ResubscribeProperties, callback func() error) error {
	result, err := c.MonitorsResubscribe(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MonitorsResubscribe: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after MonitorsResubscribe: %+v", err)
	}

	return nil
}
