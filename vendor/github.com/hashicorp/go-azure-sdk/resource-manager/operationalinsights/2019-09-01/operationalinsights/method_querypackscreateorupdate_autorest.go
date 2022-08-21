package operationalinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryPacksCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *LogAnalyticsQueryPack
}

// QueryPacksCreateOrUpdate ...
func (c OperationalInsightsClient) QueryPacksCreateOrUpdate(ctx context.Context, id QueryPackId, input LogAnalyticsQueryPack) (result QueryPacksCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForQueryPacksCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueryPacksCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueryPacksCreateOrUpdate prepares the QueryPacksCreateOrUpdate request.
func (c OperationalInsightsClient) preparerForQueryPacksCreateOrUpdate(ctx context.Context, id QueryPackId, input LogAnalyticsQueryPack) (*http.Request, error) {
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

// responderForQueryPacksCreateOrUpdate handles the response to the QueryPacksCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c OperationalInsightsClient) responderForQueryPacksCreateOrUpdate(resp *http.Response) (result QueryPacksCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
