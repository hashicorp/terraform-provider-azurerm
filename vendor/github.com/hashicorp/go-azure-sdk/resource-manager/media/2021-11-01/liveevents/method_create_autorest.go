package liveevents

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

type CreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

type CreateOperationOptions struct {
	AutoStart *bool
}

func DefaultCreateOperationOptions() CreateOperationOptions {
	return CreateOperationOptions{}
}

func (o CreateOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o CreateOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.AutoStart != nil {
		out["autoStart"] = *o.AutoStart
	}

	return out
}

// Create ...
func (c LiveEventsClient) Create(ctx context.Context, id LiveEventId, input LiveEvent, options CreateOperationOptions) (result CreateOperationResponse, err error) {
	req, err := c.preparerForCreate(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "liveevents.LiveEventsClient", "Create", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "liveevents.LiveEventsClient", "Create", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CreateThenPoll performs Create then polls until it's completed
func (c LiveEventsClient) CreateThenPoll(ctx context.Context, id LiveEventId, input LiveEvent, options CreateOperationOptions) error {
	result, err := c.Create(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing Create: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Create: %+v", err)
	}

	return nil
}

// preparerForCreate prepares the Create request.
func (c LiveEventsClient) preparerForCreate(ctx context.Context, id LiveEventId, input LiveEvent, options CreateOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCreate sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (c LiveEventsClient) senderForCreate(ctx context.Context, req *http.Request) (future CreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
