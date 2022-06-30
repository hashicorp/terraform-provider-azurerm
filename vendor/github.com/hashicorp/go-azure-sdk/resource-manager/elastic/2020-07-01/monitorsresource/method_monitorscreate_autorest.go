package monitorsresource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MonitorsCreate ...
func (c MonitorsResourceClient) MonitorsCreate(ctx context.Context, id MonitorId, input ElasticMonitorResource) (result MonitorsCreateOperationResponse, err error) {
	req, err := c.preparerForMonitorsCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMonitorsCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MonitorsCreateThenPoll performs MonitorsCreate then polls until it's completed
func (c MonitorsResourceClient) MonitorsCreateThenPoll(ctx context.Context, id MonitorId, input ElasticMonitorResource) error {
	result, err := c.MonitorsCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MonitorsCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MonitorsCreate: %+v", err)
	}

	return nil
}

// preparerForMonitorsCreate prepares the MonitorsCreate request.
func (c MonitorsResourceClient) preparerForMonitorsCreate(ctx context.Context, id MonitorId, input ElasticMonitorResource) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForMonitorsCreate sends the MonitorsCreate request. The method will close the
// http.Response Body if it receives an error.
func (c MonitorsResourceClient) senderForMonitorsCreate(ctx context.Context, req *http.Request) (future MonitorsCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
