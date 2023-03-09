package accounts

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MediaServicesOperationStatusesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *MediaServiceOperationStatus
}

// MediaServicesOperationStatusesGet ...
func (c AccountsClient) MediaServicesOperationStatusesGet(ctx context.Context, id MediaServicesOperationStatusId) (result MediaServicesOperationStatusesGetOperationResponse, err error) {
	req, err := c.preparerForMediaServicesOperationStatusesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaServicesOperationStatusesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaServicesOperationStatusesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMediaServicesOperationStatusesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaServicesOperationStatusesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMediaServicesOperationStatusesGet prepares the MediaServicesOperationStatusesGet request.
func (c AccountsClient) preparerForMediaServicesOperationStatusesGet(ctx context.Context, id MediaServicesOperationStatusId) (*http.Request, error) {
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

// responderForMediaServicesOperationStatusesGet handles the response to the MediaServicesOperationStatusesGet request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForMediaServicesOperationStatusesGet(resp *http.Response) (result MediaServicesOperationStatusesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
