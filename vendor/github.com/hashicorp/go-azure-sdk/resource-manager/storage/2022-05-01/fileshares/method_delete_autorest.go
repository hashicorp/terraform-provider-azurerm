package fileshares

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteOperationResponse struct {
	HttpResponse *http.Response
}

type DeleteOperationOptions struct {
	Include     *string
	XMsSnapshot *string
}

func DefaultDeleteOperationOptions() DeleteOperationOptions {
	return DeleteOperationOptions{}
}

func (o DeleteOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.XMsSnapshot != nil {
		out["x-ms-snapshot"] = *o.XMsSnapshot
	}

	return out
}

func (o DeleteOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Include != nil {
		out["$include"] = *o.Include
	}

	return out
}

// Delete ...
func (c FileSharesClient) Delete(ctx context.Context, id ShareId, options DeleteOperationOptions) (result DeleteOperationResponse, err error) {
	req, err := c.preparerForDelete(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "fileshares.FileSharesClient", "Delete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "fileshares.FileSharesClient", "Delete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "fileshares.FileSharesClient", "Delete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDelete prepares the Delete request.
func (c FileSharesClient) preparerForDelete(ctx context.Context, id ShareId, options DeleteOperationOptions) (*http.Request, error) {
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

// responderForDelete handles the response to the Delete request. The method always
// closes the http.Response Body.
func (c FileSharesClient) responderForDelete(resp *http.Response) (result DeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
