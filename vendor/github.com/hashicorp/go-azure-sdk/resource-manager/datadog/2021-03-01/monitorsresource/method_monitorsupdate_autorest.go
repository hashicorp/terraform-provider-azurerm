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

type MonitorsUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MonitorsUpdate ...
func (c MonitorsResourceClient) MonitorsUpdate(ctx context.Context, id MonitorId, input DatadogMonitorResourceUpdateParameters) (result MonitorsUpdateOperationResponse, err error) {
	req, err := c.preparerForMonitorsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMonitorsUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MonitorsUpdateThenPoll performs MonitorsUpdate then polls until it's completed
func (c MonitorsResourceClient) MonitorsUpdateThenPoll(ctx context.Context, id MonitorId, input DatadogMonitorResourceUpdateParameters) error {
	result, err := c.MonitorsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MonitorsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MonitorsUpdate: %+v", err)
	}

	return nil
}

// preparerForMonitorsUpdate prepares the MonitorsUpdate request.
func (c MonitorsResourceClient) preparerForMonitorsUpdate(ctx context.Context, id MonitorId, input DatadogMonitorResourceUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForMonitorsUpdate sends the MonitorsUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c MonitorsResourceClient) senderForMonitorsUpdate(ctx context.Context, req *http.Request) (future MonitorsUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
