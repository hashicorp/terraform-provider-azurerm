package accounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MediaservicesCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *MediaService
}

// MediaservicesCreateOrUpdate ...
func (c AccountsClient) MediaservicesCreateOrUpdate(ctx context.Context, id ProviderMediaServiceId, input MediaService) (result MediaservicesCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForMediaservicesCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMediaservicesCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMediaservicesCreateOrUpdate prepares the MediaservicesCreateOrUpdate request.
func (c AccountsClient) preparerForMediaservicesCreateOrUpdate(ctx context.Context, id ProviderMediaServiceId, input MediaService) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMediaservicesCreateOrUpdate handles the response to the MediaservicesCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForMediaservicesCreateOrUpdate(resp *http.Response) (result MediaservicesCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
