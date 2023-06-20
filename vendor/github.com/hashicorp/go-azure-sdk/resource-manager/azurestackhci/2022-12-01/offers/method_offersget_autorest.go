package offers

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OffersGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Offer
}

type OffersGetOperationOptions struct {
	Expand *string
}

func DefaultOffersGetOperationOptions() OffersGetOperationOptions {
	return OffersGetOperationOptions{}
}

func (o OffersGetOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o OffersGetOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Expand != nil {
		out["$expand"] = *o.Expand
	}

	return out
}

// OffersGet ...
func (c OffersClient) OffersGet(ctx context.Context, id OfferId, options OffersGetOperationOptions) (result OffersGetOperationResponse, err error) {
	req, err := c.preparerForOffersGet(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForOffersGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "offers.OffersClient", "OffersGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForOffersGet prepares the OffersGet request.
func (c OffersClient) preparerForOffersGet(ctx context.Context, id OfferId, options OffersGetOperationOptions) (*http.Request, error) {
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

// responderForOffersGet handles the response to the OffersGet request. The method always
// closes the http.Response Body.
func (c OffersClient) responderForOffersGet(resp *http.Response) (result OffersGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
