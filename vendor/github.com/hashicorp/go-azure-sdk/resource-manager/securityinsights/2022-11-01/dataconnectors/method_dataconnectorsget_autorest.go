package dataconnectors

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataConnectorsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *DataConnector
}

// DataConnectorsGet ...
func (c DataConnectorsClient) DataConnectorsGet(ctx context.Context, id DataConnectorId) (result DataConnectorsGetOperationResponse, err error) {
	req, err := c.preparerForDataConnectorsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dataconnectors.DataConnectorsClient", "DataConnectorsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dataconnectors.DataConnectorsClient", "DataConnectorsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDataConnectorsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dataconnectors.DataConnectorsClient", "DataConnectorsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDataConnectorsGet prepares the DataConnectorsGet request.
func (c DataConnectorsClient) preparerForDataConnectorsGet(ctx context.Context, id DataConnectorId) (*http.Request, error) {
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

// responderForDataConnectorsGet handles the response to the DataConnectorsGet request. The method always
// closes the http.Response Body.
func (c DataConnectorsClient) responderForDataConnectorsGet(resp *http.Response) (result DataConnectorsGetOperationResponse, err error) {
	var respObj json.RawMessage
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&respObj),
		autorest.ByClosing())
	result.HttpResponse = resp
	if err != nil {
		return
	}
	model, err := unmarshalDataConnectorImplementation(respObj)
	if err != nil {
		return
	}
	result.Model = &model
	return
}
