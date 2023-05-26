package dataconnectors

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataConnectorsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// DataConnectorsDelete ...
func (c DataConnectorsClient) DataConnectorsDelete(ctx context.Context, id DataConnectorId) (result DataConnectorsDeleteOperationResponse, err error) {
	req, err := c.preparerForDataConnectorsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dataconnectors.DataConnectorsClient", "DataConnectorsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dataconnectors.DataConnectorsClient", "DataConnectorsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDataConnectorsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dataconnectors.DataConnectorsClient", "DataConnectorsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDataConnectorsDelete prepares the DataConnectorsDelete request.
func (c DataConnectorsClient) preparerForDataConnectorsDelete(ctx context.Context, id DataConnectorId) (*http.Request, error) {
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

// responderForDataConnectorsDelete handles the response to the DataConnectorsDelete request. The method always
// closes the http.Response Body.
func (c DataConnectorsClient) responderForDataConnectorsDelete(resp *http.Response) (result DataConnectorsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
