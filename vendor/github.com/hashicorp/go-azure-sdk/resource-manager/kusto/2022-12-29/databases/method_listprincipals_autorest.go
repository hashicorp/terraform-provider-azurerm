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

type ListPrincipalsOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatabasePrincipalListResult
}

// ListPrincipals ...
func (c DatabasesClient) ListPrincipals(ctx context.Context, id DatabaseId) (result ListPrincipalsOperationResponse, err error) {
	req, err := c.preparerForListPrincipals(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ListPrincipals", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ListPrincipals", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListPrincipals(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "databases.DatabasesClient", "ListPrincipals", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListPrincipals prepares the ListPrincipals request.
func (c DatabasesClient) preparerForListPrincipals(ctx context.Context, id DatabaseId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listPrincipals", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListPrincipals handles the response to the ListPrincipals request. The method always
// closes the http.Response Body.
func (c DatabasesClient) responderForListPrincipals(resp *http.Response) (result ListPrincipalsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
