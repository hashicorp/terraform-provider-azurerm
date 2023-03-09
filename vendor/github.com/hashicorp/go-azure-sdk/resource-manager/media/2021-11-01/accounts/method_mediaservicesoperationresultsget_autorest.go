package accounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MediaServicesOperationResultsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *MediaService
}

// MediaServicesOperationResultsGet ...
func (c AccountsClient) MediaServicesOperationResultsGet(ctx context.Context, id MediaServicesOperationResultId) (result MediaServicesOperationResultsGetOperationResponse, err error) {
	req, err := c.preparerForMediaServicesOperationResultsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaServicesOperationResultsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaServicesOperationResultsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMediaServicesOperationResultsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaServicesOperationResultsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMediaServicesOperationResultsGet prepares the MediaServicesOperationResultsGet request.
func (c AccountsClient) preparerForMediaServicesOperationResultsGet(ctx context.Context, id MediaServicesOperationResultId) (*http.Request, error) {
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

// responderForMediaServicesOperationResultsGet handles the response to the MediaServicesOperationResultsGet request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForMediaServicesOperationResultsGet(resp *http.Response) (result MediaServicesOperationResultsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
