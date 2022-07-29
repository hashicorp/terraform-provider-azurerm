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

type GrafanaCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// GrafanaCreate ...
func (c GrafanaResourceClient) GrafanaCreate(ctx context.Context, id GrafanaId, input ManagedGrafana) (result GrafanaCreateOperationResponse, err error) {
	req, err := c.preparerForGrafanaCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForGrafanaCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// GrafanaCreateThenPoll performs GrafanaCreate then polls until it's completed
func (c GrafanaResourceClient) GrafanaCreateThenPoll(ctx context.Context, id GrafanaId, input ManagedGrafana) error {
	result, err := c.GrafanaCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing GrafanaCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after GrafanaCreate: %+v", err)
	}

	return nil
}

// preparerForGrafanaCreate prepares the GrafanaCreate request.
func (c GrafanaResourceClient) preparerForGrafanaCreate(ctx context.Context, id GrafanaId, input ManagedGrafana) (*http.Request, error) {
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

// senderForGrafanaCreate sends the GrafanaCreate request. The method will close the
// http.Response Body if it receives an error.
func (c GrafanaResourceClient) senderForGrafanaCreate(ctx context.Context, req *http.Request) (future GrafanaCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
