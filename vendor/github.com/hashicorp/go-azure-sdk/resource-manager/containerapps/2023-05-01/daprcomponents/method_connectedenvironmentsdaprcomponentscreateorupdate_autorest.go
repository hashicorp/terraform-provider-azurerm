package daprcomponents

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedEnvironmentsDaprComponentsCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *DaprComponent
}

// ConnectedEnvironmentsDaprComponentsCreateOrUpdate ...
func (c DaprComponentsClient) ConnectedEnvironmentsDaprComponentsCreateOrUpdate(ctx context.Context, id ConnectedEnvironmentDaprComponentId, input DaprComponent) (result ConnectedEnvironmentsDaprComponentsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForConnectedEnvironmentsDaprComponentsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "daprcomponents.DaprComponentsClient", "ConnectedEnvironmentsDaprComponentsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "daprcomponents.DaprComponentsClient", "ConnectedEnvironmentsDaprComponentsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectedEnvironmentsDaprComponentsCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "daprcomponents.DaprComponentsClient", "ConnectedEnvironmentsDaprComponentsCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectedEnvironmentsDaprComponentsCreateOrUpdate prepares the ConnectedEnvironmentsDaprComponentsCreateOrUpdate request.
func (c DaprComponentsClient) preparerForConnectedEnvironmentsDaprComponentsCreateOrUpdate(ctx context.Context, id ConnectedEnvironmentDaprComponentId, input DaprComponent) (*http.Request, error) {
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

// responderForConnectedEnvironmentsDaprComponentsCreateOrUpdate handles the response to the ConnectedEnvironmentsDaprComponentsCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c DaprComponentsClient) responderForConnectedEnvironmentsDaprComponentsCreateOrUpdate(resp *http.Response) (result ConnectedEnvironmentsDaprComponentsCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
