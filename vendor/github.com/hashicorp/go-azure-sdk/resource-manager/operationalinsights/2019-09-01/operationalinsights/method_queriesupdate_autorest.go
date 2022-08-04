package operationalinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueriesUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *LogAnalyticsQueryPackQuery
}

// QueriesUpdate ...
func (c OperationalInsightsClient) QueriesUpdate(ctx context.Context, id QueriesId, input LogAnalyticsQueryPackQuery) (result QueriesUpdateOperationResponse, err error) {
	req, err := c.preparerForQueriesUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueriesUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueriesUpdate prepares the QueriesUpdate request.
func (c OperationalInsightsClient) preparerForQueriesUpdate(ctx context.Context, id QueriesId, input LogAnalyticsQueryPackQuery) (*http.Request, error) {
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

// responderForQueriesUpdate handles the response to the QueriesUpdate request. The method always
// closes the http.Response Body.
func (c OperationalInsightsClient) responderForQueriesUpdate(resp *http.Response) (result QueriesUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
