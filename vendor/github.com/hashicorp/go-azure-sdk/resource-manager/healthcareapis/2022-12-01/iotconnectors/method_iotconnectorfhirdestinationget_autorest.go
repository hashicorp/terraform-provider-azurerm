package iotconnectors

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IotConnectorFhirDestinationGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *IotFhirDestination
}

// IotConnectorFhirDestinationGet ...
func (c IotConnectorsClient) IotConnectorFhirDestinationGet(ctx context.Context, id FhirDestinationId) (result IotConnectorFhirDestinationGetOperationResponse, err error) {
	req, err := c.preparerForIotConnectorFhirDestinationGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotconnectors.IotConnectorsClient", "IotConnectorFhirDestinationGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotconnectors.IotConnectorsClient", "IotConnectorFhirDestinationGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForIotConnectorFhirDestinationGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotconnectors.IotConnectorsClient", "IotConnectorFhirDestinationGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForIotConnectorFhirDestinationGet prepares the IotConnectorFhirDestinationGet request.
func (c IotConnectorsClient) preparerForIotConnectorFhirDestinationGet(ctx context.Context, id FhirDestinationId) (*http.Request, error) {
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

// responderForIotConnectorFhirDestinationGet handles the response to the IotConnectorFhirDestinationGet request. The method always
// closes the http.Response Body.
func (c IotConnectorsClient) responderForIotConnectorFhirDestinationGet(resp *http.Response) (result IotConnectorFhirDestinationGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
