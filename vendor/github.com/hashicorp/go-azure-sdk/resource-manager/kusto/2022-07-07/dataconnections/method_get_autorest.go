package dataconnections

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetOperationResponse struct {
	HttpResponse *http.Response
	Model        *DataConnection
}

// Get ...
func (c DataConnectionsClient) Get(ctx context.Context, id DataConnectionId) (result GetOperationResponse, err error) {
	req, err := c.preparerForGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dataconnections.DataConnectionsClient", "Get", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dataconnections.DataConnectionsClient", "Get", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dataconnections.DataConnectionsClient", "Get", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGet prepares the Get request.
func (c DataConnectionsClient) preparerForGet(ctx context.Context, id DataConnectionId) (*http.Request, error) {
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

// responderForGet handles the response to the Get request. The method always
// closes the http.Response Body.
func (c DataConnectionsClient) responderForGet(resp *http.Response) (result GetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("reading response body for DataConnection: %+v", err)
	}
	model, err := unmarshalDataConnectionImplementation(b)
	if err != nil {
		return
	}
	result.Model = &model
	return
}
