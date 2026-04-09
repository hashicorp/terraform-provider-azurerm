package subscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionAcceptOwnershipStatusOperationResponse struct {
	HttpResponse *http.Response
	Model        *AcceptOwnershipStatusResponse
}

// SubscriptionAcceptOwnershipStatus ...
func (c SubscriptionsClient) SubscriptionAcceptOwnershipStatus(ctx context.Context, id ProviderSubscriptionId) (result SubscriptionAcceptOwnershipStatusOperationResponse, err error) {
	req, err := c.preparerForSubscriptionAcceptOwnershipStatus(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionAcceptOwnershipStatus", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionAcceptOwnershipStatus", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSubscriptionAcceptOwnershipStatus(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionAcceptOwnershipStatus", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSubscriptionAcceptOwnershipStatus prepares the SubscriptionAcceptOwnershipStatus request.
func (c SubscriptionsClient) preparerForSubscriptionAcceptOwnershipStatus(ctx context.Context, id ProviderSubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/acceptOwnershipStatus", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSubscriptionAcceptOwnershipStatus handles the response to the SubscriptionAcceptOwnershipStatus request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForSubscriptionAcceptOwnershipStatus(resp *http.Response) (result SubscriptionAcceptOwnershipStatusOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
