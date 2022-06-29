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

type AccountsRegenerateKeyOperationResponse struct {
	HttpResponse *http.Response
	Model        *ApiKeys
}

// AccountsRegenerateKey ...
func (c CognitiveServicesAccountsClient) AccountsRegenerateKey(ctx context.Context, id AccountId, input RegenerateKeyParameters) (result AccountsRegenerateKeyOperationResponse, err error) {
	req, err := c.preparerForAccountsRegenerateKey(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsRegenerateKey", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsRegenerateKey", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccountsRegenerateKey(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsRegenerateKey", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccountsRegenerateKey prepares the AccountsRegenerateKey request.
func (c CognitiveServicesAccountsClient) preparerForAccountsRegenerateKey(ctx context.Context, id AccountId, input RegenerateKeyParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/regenerateKey", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAccountsRegenerateKey handles the response to the AccountsRegenerateKey request. The method always
// closes the http.Response Body.
func (c CognitiveServicesAccountsClient) responderForAccountsRegenerateKey(resp *http.Response) (result AccountsRegenerateKeyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
