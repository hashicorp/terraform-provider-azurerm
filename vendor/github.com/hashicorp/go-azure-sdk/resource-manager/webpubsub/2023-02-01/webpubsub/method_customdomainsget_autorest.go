package webpubsub

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDomainsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *CustomDomain
}

// CustomDomainsGet ...
func (c WebPubSubClient) CustomDomainsGet(ctx context.Context, id CustomDomainId) (result CustomDomainsGetOperationResponse, err error) {
	req, err := c.preparerForCustomDomainsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "CustomDomainsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "CustomDomainsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCustomDomainsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "CustomDomainsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCustomDomainsGet prepares the CustomDomainsGet request.
func (c WebPubSubClient) preparerForCustomDomainsGet(ctx context.Context, id CustomDomainId) (*http.Request, error) {
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

// responderForCustomDomainsGet handles the response to the CustomDomainsGet request. The method always
// closes the http.Response Body.
func (c WebPubSubClient) responderForCustomDomainsGet(resp *http.Response) (result CustomDomainsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
