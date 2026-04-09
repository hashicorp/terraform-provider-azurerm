package subscriptions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AliasGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *SubscriptionAliasResponse
}

// AliasGet ...
func (c SubscriptionsClient) AliasGet(ctx context.Context, id AliasId) (result AliasGetOperationResponse, err error) {
	req, err := c.preparerForAliasGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAliasGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAliasGet prepares the AliasGet request.
func (c SubscriptionsClient) preparerForAliasGet(ctx context.Context, id AliasId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAliasGet handles the response to the AliasGet request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForAliasGet(resp *http.Response) (result AliasGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
