package skuses

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkusGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Sku
}

type SkusGetOperationOptions struct {
	Expand *string
}

func DefaultSkusGetOperationOptions() SkusGetOperationOptions {
	return SkusGetOperationOptions{}
}

func (o SkusGetOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o SkusGetOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Expand != nil {
		out["$expand"] = *o.Expand
	}

	return out
}

// SkusGet ...
func (c SkusesClient) SkusGet(ctx context.Context, id SkuId, options SkusGetOperationOptions) (result SkusGetOperationResponse, err error) {
	req, err := c.preparerForSkusGet(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "skuses.SkusesClient", "SkusGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "skuses.SkusesClient", "SkusGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSkusGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "skuses.SkusesClient", "SkusGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSkusGet prepares the SkusGet request.
func (c SkusesClient) preparerForSkusGet(ctx context.Context, id SkuId, options SkusGetOperationOptions) (*http.Request, error) {
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

// responderForSkusGet handles the response to the SkusGet request. The method always
// closes the http.Response Body.
func (c SkusesClient) responderForSkusGet(resp *http.Response) (result SkusGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
