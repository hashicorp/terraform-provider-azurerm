package cognitiveservicesaccounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedAccountsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Account
}

// DeletedAccountsGet ...
func (c CognitiveServicesAccountsClient) DeletedAccountsGet(ctx context.Context, id DeletedAccountId) (result DeletedAccountsGetOperationResponse, err error) {
	req, err := c.preparerForDeletedAccountsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "DeletedAccountsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "DeletedAccountsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeletedAccountsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "DeletedAccountsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeletedAccountsGet prepares the DeletedAccountsGet request.
func (c CognitiveServicesAccountsClient) preparerForDeletedAccountsGet(ctx context.Context, id DeletedAccountId) (*http.Request, error) {
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

// responderForDeletedAccountsGet handles the response to the DeletedAccountsGet request. The method always
// closes the http.Response Body.
func (c CognitiveServicesAccountsClient) responderForDeletedAccountsGet(resp *http.Response) (result DeletedAccountsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
