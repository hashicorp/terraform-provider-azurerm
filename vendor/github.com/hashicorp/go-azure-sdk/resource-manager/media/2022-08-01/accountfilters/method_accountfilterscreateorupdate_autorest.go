package accountfilters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountFiltersCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *AccountFilter
}

// AccountFiltersCreateOrUpdate ...
func (c AccountFiltersClient) AccountFiltersCreateOrUpdate(ctx context.Context, id AccountFilterId, input AccountFilter) (result AccountFiltersCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForAccountFiltersCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccountFiltersCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccountFiltersCreateOrUpdate prepares the AccountFiltersCreateOrUpdate request.
func (c AccountFiltersClient) preparerForAccountFiltersCreateOrUpdate(ctx context.Context, id AccountFilterId, input AccountFilter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAccountFiltersCreateOrUpdate handles the response to the AccountFiltersCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c AccountFiltersClient) responderForAccountFiltersCreateOrUpdate(resp *http.Response) (result AccountFiltersCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
