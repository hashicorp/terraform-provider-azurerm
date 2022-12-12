package accounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MediaservicesDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// MediaservicesDelete ...
func (c AccountsClient) MediaservicesDelete(ctx context.Context, id ProviderMediaServiceId) (result MediaservicesDeleteOperationResponse, err error) {
	req, err := c.preparerForMediaservicesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMediaservicesDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMediaservicesDelete prepares the MediaservicesDelete request.
func (c AccountsClient) preparerForMediaservicesDelete(ctx context.Context, id ProviderMediaServiceId) (*http.Request, error) {
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

// responderForMediaservicesDelete handles the response to the MediaservicesDelete request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForMediaservicesDelete(resp *http.Response) (result MediaservicesDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
