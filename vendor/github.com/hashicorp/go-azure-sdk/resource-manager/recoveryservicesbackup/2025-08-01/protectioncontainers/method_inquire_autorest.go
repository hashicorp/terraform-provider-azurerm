package protectioncontainers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InquireOperationResponse struct {
	HttpResponse *http.Response
}

type InquireOperationOptions struct {
	Filter *string
}

func DefaultInquireOperationOptions() InquireOperationOptions {
	return InquireOperationOptions{}
}

func (o InquireOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o InquireOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// Inquire ...
func (c ProtectionContainersClient) Inquire(ctx context.Context, id ProtectionContainerId, options InquireOperationOptions) (result InquireOperationResponse, err error) {
	req, err := c.preparerForInquire(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "protectioncontainers.ProtectionContainersClient", "Inquire", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "protectioncontainers.ProtectionContainersClient", "Inquire", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForInquire(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "protectioncontainers.ProtectionContainersClient", "Inquire", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForInquire prepares the Inquire request.
func (c ProtectionContainersClient) preparerForInquire(ctx context.Context, id ProtectionContainerId, options InquireOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/inquire", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForInquire handles the response to the Inquire request. The method always
// closes the http.Response Body.
func (c ProtectionContainersClient) responderForInquire(resp *http.Response) (result InquireOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
