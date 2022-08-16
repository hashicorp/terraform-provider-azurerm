package querypackqueries

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueriesPutOperationResponse struct {
	HttpResponse *http.Response
	Model        *LogAnalyticsQueryPackQuery
}

// QueriesPut ...
func (c QueryPackQueriesClient) QueriesPut(ctx context.Context, id QueriesId, input LogAnalyticsQueryPackQuery) (result QueriesPutOperationResponse, err error) {
	req, err := c.preparerForQueriesPut(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "querypackqueries.QueryPackQueriesClient", "QueriesPut", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "querypackqueries.QueryPackQueriesClient", "QueriesPut", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueriesPut(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "querypackqueries.QueryPackQueriesClient", "QueriesPut", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueriesPut prepares the QueriesPut request.
func (c QueryPackQueriesClient) preparerForQueriesPut(ctx context.Context, id QueriesId, input LogAnalyticsQueryPackQuery) (*http.Request, error) {
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

// responderForQueriesPut handles the response to the QueriesPut request. The method always
// closes the http.Response Body.
func (c QueryPackQueriesClient) responderForQueriesPut(resp *http.Response) (result QueriesPutOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
