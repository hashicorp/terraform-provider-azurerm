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

type MonitorsLinkSaaSOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *NewRelicMonitorResource
}

// MonitorsLinkSaaS ...
func (c NewRelicMonitorResourcesClient) MonitorsLinkSaaS(ctx context.Context, id MonitorId, input SaaSData) (result MonitorsLinkSaaSOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/linkSaaS", id.ID()),
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

// MonitorsLinkSaaSThenPoll performs MonitorsLinkSaaS then polls until it's completed
func (c NewRelicMonitorResourcesClient) MonitorsLinkSaaSThenPoll(ctx context.Context, id MonitorId, input SaaSData) error {
	return c.MonitorsLinkSaaSCallbackThenPoll(ctx, id, input, nil)
}

// MonitorsLinkSaaSCallbackThenPoll performs MonitorsLinkSaaS, runs the optional callback function, then polls until it's completed
func (c NewRelicMonitorResourcesClient) MonitorsLinkSaaSCallbackThenPoll(ctx context.Context, id MonitorId, input SaaSData, callback func() error) error {
	result, err := c.MonitorsLinkSaaS(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MonitorsLinkSaaS: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after MonitorsLinkSaaS: %+v", err)
	}

	return nil
}
