package daprcomponents

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedEnvironmentsDaprComponentsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *DaprComponent
}

// ConnectedEnvironmentsDaprComponentsGet ...
func (c DaprComponentsClient) ConnectedEnvironmentsDaprComponentsGet(ctx context.Context, id ConnectedEnvironmentDaprComponentId) (result ConnectedEnvironmentsDaprComponentsGetOperationResponse, err error) {
	req, err := c.preparerForConnectedEnvironmentsDaprComponentsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "daprcomponents.DaprComponentsClient", "ConnectedEnvironmentsDaprComponentsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "daprcomponents.DaprComponentsClient", "ConnectedEnvironmentsDaprComponentsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectedEnvironmentsDaprComponentsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "daprcomponents.DaprComponentsClient", "ConnectedEnvironmentsDaprComponentsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectedEnvironmentsDaprComponentsGet prepares the ConnectedEnvironmentsDaprComponentsGet request.
func (c DaprComponentsClient) preparerForConnectedEnvironmentsDaprComponentsGet(ctx context.Context, id ConnectedEnvironmentDaprComponentId) (*http.Request, error) {
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

// responderForConnectedEnvironmentsDaprComponentsGet handles the response to the ConnectedEnvironmentsDaprComponentsGet request. The method always
// closes the http.Response Body.
func (c DaprComponentsClient) responderForConnectedEnvironmentsDaprComponentsGet(resp *http.Response) (result ConnectedEnvironmentsDaprComponentsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
