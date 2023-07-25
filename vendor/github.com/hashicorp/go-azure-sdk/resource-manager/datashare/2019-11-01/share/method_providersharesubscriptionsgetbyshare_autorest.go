package share

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderShareSubscriptionsGetByShareOperationResponse struct {
	HttpResponse *http.Response
	Model        *ProviderShareSubscription
}

// ProviderShareSubscriptionsGetByShare ...
func (c ShareClient) ProviderShareSubscriptionsGetByShare(ctx context.Context, id ProviderShareSubscriptionId) (result ProviderShareSubscriptionsGetByShareOperationResponse, err error) {
	req, err := c.preparerForProviderShareSubscriptionsGetByShare(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ProviderShareSubscriptionsGetByShare", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ProviderShareSubscriptionsGetByShare", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForProviderShareSubscriptionsGetByShare(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ProviderShareSubscriptionsGetByShare", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForProviderShareSubscriptionsGetByShare prepares the ProviderShareSubscriptionsGetByShare request.
func (c ShareClient) preparerForProviderShareSubscriptionsGetByShare(ctx context.Context, id ProviderShareSubscriptionId) (*http.Request, error) {
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

// responderForProviderShareSubscriptionsGetByShare handles the response to the ProviderShareSubscriptionsGetByShare request. The method always
// closes the http.Response Body.
func (c ShareClient) responderForProviderShareSubscriptionsGetByShare(resp *http.Response) (result ProviderShareSubscriptionsGetByShareOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
