package daprcomponents

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedEnvironmentsDaprComponentsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// ConnectedEnvironmentsDaprComponentsDelete ...
func (c DaprComponentsClient) ConnectedEnvironmentsDaprComponentsDelete(ctx context.Context, id ConnectedEnvironmentDaprComponentId) (result ConnectedEnvironmentsDaprComponentsDeleteOperationResponse, err error) {
	req, err := c.preparerForConnectedEnvironmentsDaprComponentsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "daprcomponents.DaprComponentsClient", "ConnectedEnvironmentsDaprComponentsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "daprcomponents.DaprComponentsClient", "ConnectedEnvironmentsDaprComponentsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectedEnvironmentsDaprComponentsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "daprcomponents.DaprComponentsClient", "ConnectedEnvironmentsDaprComponentsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectedEnvironmentsDaprComponentsDelete prepares the ConnectedEnvironmentsDaprComponentsDelete request.
func (c DaprComponentsClient) preparerForConnectedEnvironmentsDaprComponentsDelete(ctx context.Context, id ConnectedEnvironmentDaprComponentId) (*http.Request, error) {
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

// responderForConnectedEnvironmentsDaprComponentsDelete handles the response to the ConnectedEnvironmentsDaprComponentsDelete request. The method always
// closes the http.Response Body.
func (c DaprComponentsClient) responderForConnectedEnvironmentsDaprComponentsDelete(resp *http.Response) (result ConnectedEnvironmentsDaprComponentsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
