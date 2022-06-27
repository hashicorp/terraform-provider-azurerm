package cognitiveservicesaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountsListSkusOperationResponse struct {
	HttpResponse *http.Response
	Model        *AccountSkuListResult
}

// AccountsListSkus ...
func (c CognitiveServicesAccountsClient) AccountsListSkus(ctx context.Context, id AccountId) (result AccountsListSkusOperationResponse, err error) {
	req, err := c.preparerForAccountsListSkus(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListSkus", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListSkus", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccountsListSkus(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListSkus", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccountsListSkus prepares the AccountsListSkus request.
func (c CognitiveServicesAccountsClient) preparerForAccountsListSkus(ctx context.Context, id AccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/skus", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAccountsListSkus handles the response to the AccountsListSkus request. The method always
// closes the http.Response Body.
func (c CognitiveServicesAccountsClient) responderForAccountsListSkus(resp *http.Response) (result AccountsListSkusOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
