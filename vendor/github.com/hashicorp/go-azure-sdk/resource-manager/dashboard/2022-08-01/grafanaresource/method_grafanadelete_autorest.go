package grafanaresource

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

type GrafanaDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// GrafanaDelete ...
func (c GrafanaResourceClient) GrafanaDelete(ctx context.Context, id GrafanaId) (result GrafanaDeleteOperationResponse, err error) {
	req, err := c.preparerForGrafanaDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForGrafanaDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// GrafanaDeleteThenPoll performs GrafanaDelete then polls until it's completed
func (c GrafanaResourceClient) GrafanaDeleteThenPoll(ctx context.Context, id GrafanaId) error {
	result, err := c.GrafanaDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing GrafanaDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after GrafanaDelete: %+v", err)
	}

	return nil
}

// preparerForGrafanaDelete prepares the GrafanaDelete request.
func (c GrafanaResourceClient) preparerForGrafanaDelete(ctx context.Context, id GrafanaId) (*http.Request, error) {
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

// senderForGrafanaDelete sends the GrafanaDelete request. The method will close the
// http.Response Body if it receives an error.
func (c GrafanaResourceClient) senderForGrafanaDelete(ctx context.Context, req *http.Request) (future GrafanaDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
