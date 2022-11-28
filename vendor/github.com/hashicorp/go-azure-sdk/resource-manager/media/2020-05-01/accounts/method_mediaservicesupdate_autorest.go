package accounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MediaservicesUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *MediaService
}

// MediaservicesUpdate ...
func (c AccountsClient) MediaservicesUpdate(ctx context.Context, id ProviderMediaServiceId, input MediaService) (result MediaservicesUpdateOperationResponse, err error) {
	req, err := c.preparerForMediaservicesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMediaservicesUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMediaservicesUpdate prepares the MediaservicesUpdate request.
func (c AccountsClient) preparerForMediaservicesUpdate(ctx context.Context, id ProviderMediaServiceId, input MediaService) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMediaservicesUpdate handles the response to the MediaservicesUpdate request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForMediaservicesUpdate(resp *http.Response) (result MediaservicesUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
