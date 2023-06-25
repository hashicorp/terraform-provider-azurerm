package connectedclusters

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedClusterUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *ConnectedCluster
}

// ConnectedClusterUpdate ...
func (c ConnectedClustersClient) ConnectedClusterUpdate(ctx context.Context, id ConnectedClusterId, input ConnectedClusterPatch) (result ConnectedClusterUpdateOperationResponse, err error) {
	req, err := c.preparerForConnectedClusterUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForConnectedClusterUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForConnectedClusterUpdate prepares the ConnectedClusterUpdate request.
func (c ConnectedClustersClient) preparerForConnectedClusterUpdate(ctx context.Context, id ConnectedClusterId, input ConnectedClusterPatch) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForConnectedClusterUpdate handles the response to the ConnectedClusterUpdate request. The method always
// closes the http.Response Body.
func (c ConnectedClustersClient) responderForConnectedClusterUpdate(resp *http.Response) (result ConnectedClusterUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
