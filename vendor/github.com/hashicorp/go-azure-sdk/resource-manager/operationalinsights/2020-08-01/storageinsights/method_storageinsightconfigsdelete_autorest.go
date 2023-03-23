package storageinsights

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageInsightConfigsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// StorageInsightConfigsDelete ...
func (c StorageInsightsClient) StorageInsightConfigsDelete(ctx context.Context, id StorageInsightConfigId) (result StorageInsightConfigsDeleteOperationResponse, err error) {
	req, err := c.preparerForStorageInsightConfigsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageinsights.StorageInsightsClient", "StorageInsightConfigsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageinsights.StorageInsightsClient", "StorageInsightConfigsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForStorageInsightConfigsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storageinsights.StorageInsightsClient", "StorageInsightConfigsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForStorageInsightConfigsDelete prepares the StorageInsightConfigsDelete request.
func (c StorageInsightsClient) preparerForStorageInsightConfigsDelete(ctx context.Context, id StorageInsightConfigId) (*http.Request, error) {
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

// responderForStorageInsightConfigsDelete handles the response to the StorageInsightConfigsDelete request. The method always
// closes the http.Response Body.
func (c StorageInsightsClient) responderForStorageInsightConfigsDelete(resp *http.Response) (result StorageInsightConfigsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
