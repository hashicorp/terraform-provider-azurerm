package querypacks

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryPacksGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *LogAnalyticsQueryPack
}

// QueryPacksGet ...
func (c QueryPacksClient) QueryPacksGet(ctx context.Context, id QueryPackId) (result QueryPacksGetOperationResponse, err error) {
	req, err := c.preparerForQueryPacksGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "querypacks.QueryPacksClient", "QueryPacksGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "querypacks.QueryPacksClient", "QueryPacksGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForQueryPacksGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "querypacks.QueryPacksClient", "QueryPacksGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForQueryPacksGet prepares the QueryPacksGet request.
func (c QueryPacksClient) preparerForQueryPacksGet(ctx context.Context, id QueryPackId) (*http.Request, error) {
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

// responderForQueryPacksGet handles the response to the QueryPacksGet request. The method always
// closes the http.Response Body.
func (c QueryPacksClient) responderForQueryPacksGet(resp *http.Response) (result QueryPacksGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
