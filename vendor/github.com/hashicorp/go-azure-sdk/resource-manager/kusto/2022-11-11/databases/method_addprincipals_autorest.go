package databases

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AddPrincipalsOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatabasePrincipalListResult
}

// AddPrincipals ...
func (c DatabasesClient) AddPrincipals(ctx context.Context, id DatabaseId, input DatabasePrincipalListRequest) (result AddPrincipalsOperationResponse, err error) {
	req, err := c.preparerForAddPrincipals(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "AddPrincipals", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "AddPrincipals", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAddPrincipals(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "AddPrincipals", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAddPrincipals prepares the AddPrincipals request.
func (c DatabasesClient) preparerForAddPrincipals(ctx context.Context, id DatabaseId, input DatabasePrincipalListRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/addPrincipals", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAddPrincipals handles the response to the AddPrincipals request. The method always
// closes the http.Response Body.
func (c DatabasesClient) responderForAddPrincipals(resp *http.Response) (result AddPrincipalsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
