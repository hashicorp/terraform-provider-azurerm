package accounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MediaservicesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *MediaService
}

// MediaservicesGet ...
func (c AccountsClient) MediaservicesGet(ctx context.Context, id MediaServiceId) (result MediaservicesGetOperationResponse, err error) {
	req, err := c.preparerForMediaservicesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMediaservicesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMediaservicesGet prepares the MediaservicesGet request.
func (c AccountsClient) preparerForMediaservicesGet(ctx context.Context, id MediaServiceId) (*http.Request, error) {
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

// responderForMediaservicesGet handles the response to the MediaservicesGet request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForMediaservicesGet(resp *http.Response) (result MediaservicesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
