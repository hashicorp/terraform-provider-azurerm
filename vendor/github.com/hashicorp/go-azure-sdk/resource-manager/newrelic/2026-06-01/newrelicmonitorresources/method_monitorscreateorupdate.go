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

type MonitorsCreateOrUpdateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *NewRelicMonitorResource
}

// MonitorsCreateOrUpdate ...
func (c NewRelicMonitorResourcesClient) MonitorsCreateOrUpdate(ctx context.Context, id MonitorId, input NewRelicMonitorResource) (result MonitorsCreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       id.ID(),
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

// MonitorsCreateOrUpdateThenPoll performs MonitorsCreateOrUpdate then polls until it's completed
func (c NewRelicMonitorResourcesClient) MonitorsCreateOrUpdateThenPoll(ctx context.Context, id MonitorId, input NewRelicMonitorResource) error {
	return c.MonitorsCreateOrUpdateCallbackThenPoll(ctx, id, input, nil)
}

// MonitorsCreateOrUpdateCallbackThenPoll performs MonitorsCreateOrUpdate, runs the optional callback function, then polls until it's completed
func (c NewRelicMonitorResourcesClient) MonitorsCreateOrUpdateCallbackThenPoll(ctx context.Context, id MonitorId, input NewRelicMonitorResource, callback func() error) error {
	result, err := c.MonitorsCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MonitorsCreateOrUpdate: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after MonitorsCreateOrUpdate: %+v", err)
	}

	return nil
}
