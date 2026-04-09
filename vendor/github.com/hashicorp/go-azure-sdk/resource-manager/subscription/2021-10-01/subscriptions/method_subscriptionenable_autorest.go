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

type SubscriptionEnableOperationResponse struct {
	HttpResponse *http.Response
	Model        *EnabledSubscriptionId
}

// SubscriptionEnable ...
func (c SubscriptionsClient) SubscriptionEnable(ctx context.Context, id commonids.SubscriptionId) (result SubscriptionEnableOperationResponse, err error) {
	req, err := c.preparerForSubscriptionEnable(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionEnable", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionEnable", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSubscriptionEnable(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionEnable", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSubscriptionEnable prepares the SubscriptionEnable request.
func (c SubscriptionsClient) preparerForSubscriptionEnable(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Subscription/enable", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSubscriptionEnable handles the response to the SubscriptionEnable request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForSubscriptionEnable(resp *http.Response) (result SubscriptionEnableOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
