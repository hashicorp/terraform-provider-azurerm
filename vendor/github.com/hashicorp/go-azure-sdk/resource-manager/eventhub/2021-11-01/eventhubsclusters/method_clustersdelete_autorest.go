package eventhubsclusters

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

type ClustersDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ClustersDelete ...
func (c EventHubsClustersClient) ClustersDelete(ctx context.Context, id ClusterId) (result ClustersDeleteOperationResponse, err error) {
	req, err := c.preparerForClustersDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForClustersDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ClustersDeleteThenPoll performs ClustersDelete then polls until it's completed
func (c EventHubsClustersClient) ClustersDeleteThenPoll(ctx context.Context, id ClusterId) error {
	result, err := c.ClustersDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing ClustersDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ClustersDelete: %+v", err)
	}

	return nil
}

// preparerForClustersDelete prepares the ClustersDelete request.
func (c EventHubsClustersClient) preparerForClustersDelete(ctx context.Context, id ClusterId) (*http.Request, error) {
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

// senderForClustersDelete sends the ClustersDelete request. The method will close the
// http.Response Body if it receives an error.
func (c EventHubsClustersClient) senderForClustersDelete(ctx context.Context, req *http.Request) (future ClustersDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
