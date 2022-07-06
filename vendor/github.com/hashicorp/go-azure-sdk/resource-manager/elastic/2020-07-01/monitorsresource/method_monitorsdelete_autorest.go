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

type MonitorsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MonitorsDelete ...
func (c MonitorsResourceClient) MonitorsDelete(ctx context.Context, id MonitorId) (result MonitorsDeleteOperationResponse, err error) {
	req, err := c.preparerForMonitorsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMonitorsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MonitorsDeleteThenPoll performs MonitorsDelete then polls until it's completed
func (c MonitorsResourceClient) MonitorsDeleteThenPoll(ctx context.Context, id MonitorId) error {
	result, err := c.MonitorsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing MonitorsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MonitorsDelete: %+v", err)
	}

	return nil
}

// preparerForMonitorsDelete prepares the MonitorsDelete request.
func (c MonitorsResourceClient) preparerForMonitorsDelete(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForMonitorsDelete sends the MonitorsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c MonitorsResourceClient) senderForMonitorsDelete(ctx context.Context, req *http.Request) (future MonitorsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
