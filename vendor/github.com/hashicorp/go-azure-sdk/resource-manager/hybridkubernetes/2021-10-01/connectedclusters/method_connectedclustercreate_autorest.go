package connectedclusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedClusterCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ConnectedClusterCreate ...
func (c ConnectedClustersClient) ConnectedClusterCreate(ctx context.Context, id ConnectedClusterId, input ConnectedCluster) (result ConnectedClusterCreateOperationResponse, err error) {
	req, err := c.preparerForConnectedClusterCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForConnectedClusterCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "connectedclusters.ConnectedClustersClient", "ConnectedClusterCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ConnectedClusterCreateThenPoll performs ConnectedClusterCreate then polls until it's completed
func (c ConnectedClustersClient) ConnectedClusterCreateThenPoll(ctx context.Context, id ConnectedClusterId, input ConnectedCluster) error {
	result, err := c.ConnectedClusterCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ConnectedClusterCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ConnectedClusterCreate: %+v", err)
	}

	return nil
}

// preparerForConnectedClusterCreate prepares the ConnectedClusterCreate request.
func (c ConnectedClustersClient) preparerForConnectedClusterCreate(ctx context.Context, id ConnectedClusterId, input ConnectedCluster) (*http.Request, error) {
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

// senderForConnectedClusterCreate sends the ConnectedClusterCreate request. The method will close the
// http.Response Body if it receives an error.
func (c ConnectedClustersClient) senderForConnectedClusterCreate(ctx context.Context, req *http.Request) (future ConnectedClusterCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
