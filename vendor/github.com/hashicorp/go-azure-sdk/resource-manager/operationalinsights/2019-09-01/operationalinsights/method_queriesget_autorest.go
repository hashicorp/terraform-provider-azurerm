package operationalinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueriesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *LogAnalyticsQueryPackQuery
}

// QueriesGet ...
func (c OperationalInsightsClient) QueriesGet(ctx context.Context, id QueriesId) (result QueriesGetOperationResponse, err error) {
	req, err := c.preparerForQueriesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueriesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueriesGet prepares the QueriesGet request.
func (c OperationalInsightsClient) preparerForQueriesGet(ctx context.Context, id QueriesId) (*http.Request, error) {
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

// responderForQueriesGet handles the response to the QueriesGet request. The method always
// closes the http.Response Body.
func (c OperationalInsightsClient) responderForQueriesGet(resp *http.Response) (result QueriesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
