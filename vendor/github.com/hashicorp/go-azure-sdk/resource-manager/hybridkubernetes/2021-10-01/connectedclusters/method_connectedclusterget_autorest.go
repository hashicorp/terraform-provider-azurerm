package connectedclusters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedClusterGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConnectedCluster
}

// ConnectedClusterGet ...
func (c ConnectedClustersClient) ConnectedClusterGet(ctx context.Context, id ConnectedClusterId) (result ConnectedClusterGetOperationResponse, err error) {
	req, err := c.preparerForConnectedClusterGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectedClusterGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectedClusterGet prepares the ConnectedClusterGet request.
func (c ConnectedClustersClient) preparerForConnectedClusterGet(ctx context.Context, id ConnectedClusterId) (*http.Request, error) {
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

// responderForConnectedClusterGet handles the response to the ConnectedClusterGet request. The method always
// closes the http.Response Body.
func (c ConnectedClustersClient) responderForConnectedClusterGet(resp *http.Response) (result ConnectedClusterGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
