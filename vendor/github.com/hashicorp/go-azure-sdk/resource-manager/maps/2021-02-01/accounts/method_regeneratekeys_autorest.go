package accounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegenerateKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *MapsAccountKeys
}

// RegenerateKeys ...
func (c AccountsClient) RegenerateKeys(ctx context.Context, id AccountId, input MapsKeySpecification) (result RegenerateKeysOperationResponse, err error) {
	req, err := c.preparerForRegenerateKeys(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "RegenerateKeys", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "RegenerateKeys", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRegenerateKeys(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "RegenerateKeys", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRegenerateKeys prepares the RegenerateKeys request.
func (c AccountsClient) preparerForRegenerateKeys(ctx context.Context, id AccountId, input MapsKeySpecification) (*http.Request, error) {
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

// responderForRegenerateKeys handles the response to the RegenerateKeys request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForRegenerateKeys(resp *http.Response) (result RegenerateKeysOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
