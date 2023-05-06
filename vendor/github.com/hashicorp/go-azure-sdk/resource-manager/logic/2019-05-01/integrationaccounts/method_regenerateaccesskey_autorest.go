package integrationaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegenerateAccessKeyOperationResponse struct {
	HttpResponse *http.Response
	Model        *IntegrationAccount
}

// RegenerateAccessKey ...
func (c IntegrationAccountsClient) RegenerateAccessKey(ctx context.Context, id IntegrationAccountId, input RegenerateActionParameter) (result RegenerateAccessKeyOperationResponse, err error) {
	req, err := c.preparerForRegenerateAccessKey(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "integrationaccounts.IntegrationAccountsClient", "RegenerateAccessKey", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "integrationaccounts.IntegrationAccountsClient", "RegenerateAccessKey", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRegenerateAccessKey(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "integrationaccounts.IntegrationAccountsClient", "RegenerateAccessKey", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRegenerateAccessKey prepares the RegenerateAccessKey request.
func (c IntegrationAccountsClient) preparerForRegenerateAccessKey(ctx context.Context, id IntegrationAccountId, input RegenerateActionParameter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/regenerateAccessKey", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRegenerateAccessKey handles the response to the RegenerateAccessKey request. The method always
// closes the http.Response Body.
func (c IntegrationAccountsClient) responderForRegenerateAccessKey(resp *http.Response) (result RegenerateAccessKeyOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
