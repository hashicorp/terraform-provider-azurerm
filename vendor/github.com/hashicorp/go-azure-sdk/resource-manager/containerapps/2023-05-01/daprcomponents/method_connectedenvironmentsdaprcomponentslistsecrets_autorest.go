package daprcomponents

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedEnvironmentsDaprComponentsListSecretsOperationResponse struct {
	HttpResponse *http.Response
	Model        *DaprSecretsCollection
}

// ConnectedEnvironmentsDaprComponentsListSecrets ...
func (c DaprComponentsClient) ConnectedEnvironmentsDaprComponentsListSecrets(ctx context.Context, id ConnectedEnvironmentDaprComponentId) (result ConnectedEnvironmentsDaprComponentsListSecretsOperationResponse, err error) {
	req, err := c.preparerForConnectedEnvironmentsDaprComponentsListSecrets(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "daprcomponents.DaprComponentsClient", "ConnectedEnvironmentsDaprComponentsListSecrets", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "daprcomponents.DaprComponentsClient", "ConnectedEnvironmentsDaprComponentsListSecrets", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectedEnvironmentsDaprComponentsListSecrets(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "daprcomponents.DaprComponentsClient", "ConnectedEnvironmentsDaprComponentsListSecrets", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectedEnvironmentsDaprComponentsListSecrets prepares the ConnectedEnvironmentsDaprComponentsListSecrets request.
func (c DaprComponentsClient) preparerForConnectedEnvironmentsDaprComponentsListSecrets(ctx context.Context, id ConnectedEnvironmentDaprComponentId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listSecrets", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConnectedEnvironmentsDaprComponentsListSecrets handles the response to the ConnectedEnvironmentsDaprComponentsListSecrets request. The method always
// closes the http.Response Body.
func (c DaprComponentsClient) responderForConnectedEnvironmentsDaprComponentsListSecrets(resp *http.Response) (result ConnectedEnvironmentsDaprComponentsListSecretsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
