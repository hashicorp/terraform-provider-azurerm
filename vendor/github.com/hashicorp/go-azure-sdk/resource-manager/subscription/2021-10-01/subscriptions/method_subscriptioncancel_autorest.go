package subscriptions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubscriptionCancelOperationResponse struct {
	HttpResponse *http.Response
	Model        *CanceledSubscriptionId
}

// SubscriptionCancel ...
func (c SubscriptionsClient) SubscriptionCancel(ctx context.Context, id commonids.SubscriptionId) (result SubscriptionCancelOperationResponse, err error) {
	req, err := c.preparerForSubscriptionCancel(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionCancel", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionCancel", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSubscriptionCancel(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionCancel", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSubscriptionCancel prepares the SubscriptionCancel request.
func (c SubscriptionsClient) preparerForSubscriptionCancel(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Subscription/cancel", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSubscriptionCancel handles the response to the SubscriptionCancel request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForSubscriptionCancel(resp *http.Response) (result SubscriptionCancelOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
