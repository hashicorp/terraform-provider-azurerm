package accountfilters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountFiltersUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *AccountFilter
}

// AccountFiltersUpdate ...
func (c AccountFiltersClient) AccountFiltersUpdate(ctx context.Context, id AccountFilterId, input AccountFilter) (result AccountFiltersUpdateOperationResponse, err error) {
	req, err := c.preparerForAccountFiltersUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAccountFiltersUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accountfilters.AccountFiltersClient", "AccountFiltersUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAccountFiltersUpdate prepares the AccountFiltersUpdate request.
func (c AccountFiltersClient) preparerForAccountFiltersUpdate(ctx context.Context, id AccountFilterId, input AccountFilter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForAccountFiltersUpdate handles the response to the AccountFiltersUpdate request. The method always
// closes the http.Response Body.
func (c AccountFiltersClient) responderForAccountFiltersUpdate(resp *http.Response) (result AccountFiltersUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
