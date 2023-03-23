package sharedprivatelinkresources

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetOperationResponse struct {
	HttpResponse *http.Response
	Model        *SharedPrivateLinkResource
}

type GetOperationOptions struct {
	XMsClientRequestId *string
}

func DefaultGetOperationOptions() GetOperationOptions {
	return GetOperationOptions{}
}

func (o GetOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.XMsClientRequestId != nil {
		out["x-ms-client-request-id"] = *o.XMsClientRequestId
	}

	return out
}

func (o GetOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// Get ...
func (c SharedPrivateLinkResourcesClient) Get(ctx context.Context, id SharedPrivateLinkResourceId, options GetOperationOptions) (result GetOperationResponse, err error) {
	req, err := c.preparerForGet(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sharedprivatelinkresources.SharedPrivateLinkResourcesClient", "Get", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "sharedprivatelinkresources.SharedPrivateLinkResourcesClient", "Get", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sharedprivatelinkresources.SharedPrivateLinkResourcesClient", "Get", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGet prepares the Get request.
func (c SharedPrivateLinkResourcesClient) preparerForGet(ctx context.Context, id SharedPrivateLinkResourceId, options GetOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGet handles the response to the Get request. The method always
// closes the http.Response Body.
func (c SharedPrivateLinkResourcesClient) responderForGet(resp *http.Response) (result GetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
