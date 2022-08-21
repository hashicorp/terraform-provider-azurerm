package operationalinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryPacksDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// QueryPacksDelete ...
func (c OperationalInsightsClient) QueryPacksDelete(ctx context.Context, id QueryPackId) (result QueryPacksDeleteOperationResponse, err error) {
	req, err := c.preparerForQueryPacksDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueryPacksDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueryPacksDelete prepares the QueryPacksDelete request.
func (c OperationalInsightsClient) preparerForQueryPacksDelete(ctx context.Context, id QueryPackId) (*http.Request, error) {
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

// responderForQueryPacksDelete handles the response to the QueryPacksDelete request. The method always
// closes the http.Response Body.
func (c OperationalInsightsClient) responderForQueryPacksDelete(resp *http.Response) (result QueryPacksDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
