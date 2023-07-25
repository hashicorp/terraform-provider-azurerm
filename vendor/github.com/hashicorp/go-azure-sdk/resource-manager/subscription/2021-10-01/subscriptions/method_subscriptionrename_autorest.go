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

type SubscriptionRenameOperationResponse struct {
	HttpResponse *http.Response
	Model        *RenamedSubscriptionId
}

// SubscriptionRename ...
func (c SubscriptionsClient) SubscriptionRename(ctx context.Context, id commonids.SubscriptionId, input SubscriptionName) (result SubscriptionRenameOperationResponse, err error) {
	req, err := c.preparerForSubscriptionRename(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionRename", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionRename", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSubscriptionRename(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "SubscriptionRename", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSubscriptionRename prepares the SubscriptionRename request.
func (c SubscriptionsClient) preparerForSubscriptionRename(ctx context.Context, id commonids.SubscriptionId, input SubscriptionName) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Subscription/rename", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSubscriptionRename handles the response to the SubscriptionRename request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForSubscriptionRename(resp *http.Response) (result SubscriptionRenameOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
