package fileshares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LeaseOperationResponse struct {
	HttpResponse *http.Response
	Model        *LeaseShareResponse
}

type LeaseOperationOptions struct {
	XMsSnapshot *string
}

func DefaultLeaseOperationOptions() LeaseOperationOptions {
	return LeaseOperationOptions{}
}

func (o LeaseOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.XMsSnapshot != nil {
		out["x-ms-snapshot"] = *o.XMsSnapshot
	}

	return out
}

func (o LeaseOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// Lease ...
func (c FileSharesClient) Lease(ctx context.Context, id ShareId, input LeaseShareRequest, options LeaseOperationOptions) (result LeaseOperationResponse, err error) {
	req, err := c.preparerForLease(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "fileshares.FileSharesClient", "Lease", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "fileshares.FileSharesClient", "Lease", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForLease(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "fileshares.FileSharesClient", "Lease", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForLease prepares the Lease request.
func (c FileSharesClient) preparerForLease(ctx context.Context, id ShareId, input LeaseShareRequest, options LeaseOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/lease", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForLease handles the response to the Lease request. The method always
// closes the http.Response Body.
func (c FileSharesClient) responderForLease(resp *http.Response) (result LeaseOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
