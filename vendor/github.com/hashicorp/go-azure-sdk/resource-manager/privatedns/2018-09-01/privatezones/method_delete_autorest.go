package privatezones

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

type DeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

type DeleteOperationOptions struct {
	IfMatch *string
}

func DefaultDeleteOperationOptions() DeleteOperationOptions {
	return DeleteOperationOptions{}
}

func (o DeleteOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	return out
}

func (o DeleteOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// Delete ...
func (c PrivateZonesClient) Delete(ctx context.Context, id PrivateDnsZoneId, options DeleteOperationOptions) (result DeleteOperationResponse, err error) {
	req, err := c.preparerForDelete(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatezones.PrivateZonesClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatezones.PrivateZonesClient", "Delete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DeleteThenPoll performs Delete then polls until it's completed
func (c PrivateZonesClient) DeleteThenPoll(ctx context.Context, id PrivateDnsZoneId, options DeleteOperationOptions) error {
	result, err := c.Delete(ctx, id, options)
	if err != nil {
		return fmt.Errorf("performing Delete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Delete: %+v", err)
	}

	return nil
}

// preparerForDelete prepares the Delete request.
func (c PrivateZonesClient) preparerForDelete(ctx context.Context, id PrivateDnsZoneId, options DeleteOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDelete sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (c PrivateZonesClient) senderForDelete(ctx context.Context, req *http.Request) (future DeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
