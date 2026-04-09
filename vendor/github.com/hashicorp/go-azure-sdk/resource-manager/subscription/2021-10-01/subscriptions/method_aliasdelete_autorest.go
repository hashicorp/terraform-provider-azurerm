package subscriptions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AliasDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// AliasDelete ...
func (c SubscriptionsClient) AliasDelete(ctx context.Context, id AliasId) (result AliasDeleteOperationResponse, err error) {
	req, err := c.preparerForAliasDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAliasDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAliasDelete prepares the AliasDelete request.
func (c SubscriptionsClient) preparerForAliasDelete(ctx context.Context, id AliasId) (*http.Request, error) {
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

// responderForAliasDelete handles the response to the AliasDelete request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForAliasDelete(resp *http.Response) (result AliasDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
