package cognitiveservicesaccounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Account
}

// AccountsGet ...
func (c CognitiveServicesAccountsClient) AccountsGet(ctx context.Context, id AccountId) (result AccountsGetOperationResponse, err error) {
	req, err := c.preparerForAccountsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccountsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccountsGet prepares the AccountsGet request.
func (c CognitiveServicesAccountsClient) preparerForAccountsGet(ctx context.Context, id AccountId) (*http.Request, error) {
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

// responderForAccountsGet handles the response to the AccountsGet request. The method always
// closes the http.Response Body.
func (c CognitiveServicesAccountsClient) responderForAccountsGet(resp *http.Response) (result AccountsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
