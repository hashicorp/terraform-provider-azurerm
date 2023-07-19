package subscriptions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionOperationGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *SubscriptionCreationResult
}

// SubscriptionOperationGet ...
func (c SubscriptionsClient) SubscriptionOperationGet(ctx context.Context, id SubscriptionOperationId) (result SubscriptionOperationGetOperationResponse, err error) {
	req, err := c.preparerForSubscriptionOperationGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionOperationGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionOperationGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSubscriptionOperationGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionOperationGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSubscriptionOperationGet prepares the SubscriptionOperationGet request.
func (c SubscriptionsClient) preparerForSubscriptionOperationGet(ctx context.Context, id SubscriptionOperationId) (*http.Request, error) {
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

// responderForSubscriptionOperationGet handles the response to the SubscriptionOperationGet request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForSubscriptionOperationGet(resp *http.Response) (result SubscriptionOperationGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
