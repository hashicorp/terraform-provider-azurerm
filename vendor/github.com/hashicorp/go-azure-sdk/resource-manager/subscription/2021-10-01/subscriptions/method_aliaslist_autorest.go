package subscriptions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AliasListOperationResponse struct {
	HttpResponse *http.Response
	Model        *SubscriptionAliasListResult
}

// AliasList ...
func (c SubscriptionsClient) AliasList(ctx context.Context) (result AliasListOperationResponse, err error) {
	req, err := c.preparerForAliasList(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAliasList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAliasList prepares the AliasList request.
func (c SubscriptionsClient) preparerForAliasList(ctx context.Context) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath("/providers/Microsoft.Subscription/aliases"),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAliasList handles the response to the AliasList request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForAliasList(resp *http.Response) (result AliasListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
