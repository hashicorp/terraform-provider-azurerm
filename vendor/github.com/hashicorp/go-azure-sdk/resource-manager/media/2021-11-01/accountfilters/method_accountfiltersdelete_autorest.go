package accountfilters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountFiltersDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// AccountFiltersDelete ...
func (c AccountFiltersClient) AccountFiltersDelete(ctx context.Context, id AccountFilterId) (result AccountFiltersDeleteOperationResponse, err error) {
	req, err := c.preparerForAccountFiltersDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccountFiltersDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccountFiltersDelete prepares the AccountFiltersDelete request.
func (c AccountFiltersClient) preparerForAccountFiltersDelete(ctx context.Context, id AccountFilterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAccountFiltersDelete handles the response to the AccountFiltersDelete request. The method always
// closes the http.Response Body.
func (c AccountFiltersClient) responderForAccountFiltersDelete(resp *http.Response) (result AccountFiltersDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
