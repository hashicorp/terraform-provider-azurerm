package cosmosdb

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *LocationGetResult
}

// LocationsGet ...
func (c CosmosDBClient) LocationsGet(ctx context.Context, id LocationId) (result LocationsGetOperationResponse, err error) {
	req, err := c.preparerForLocationsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "LocationsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "LocationsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForLocationsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cosmosdb.CosmosDBClient", "LocationsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForLocationsGet prepares the LocationsGet request.
func (c CosmosDBClient) preparerForLocationsGet(ctx context.Context, id LocationId) (*http.Request, error) {
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

// responderForLocationsGet handles the response to the LocationsGet request. The method always
// closes the http.Response Body.
func (c CosmosDBClient) responderForLocationsGet(resp *http.Response) (result LocationsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
