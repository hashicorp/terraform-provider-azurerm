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

type ClustersUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ClustersUpdate ...
func (c EventHubsClustersClient) ClustersUpdate(ctx context.Context, id ClusterId, input Cluster) (result ClustersUpdateOperationResponse, err error) {
	req, err := c.preparerForClustersUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForClustersUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubsclusters.EventHubsClustersClient", "ClustersUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ClustersUpdateThenPoll performs ClustersUpdate then polls until it's completed
func (c EventHubsClustersClient) ClustersUpdateThenPoll(ctx context.Context, id ClusterId, input Cluster) error {
	result, err := c.ClustersUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ClustersUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ClustersUpdate: %+v", err)
	}

	return nil
}

// preparerForClustersUpdate prepares the ClustersUpdate request.
func (c EventHubsClustersClient) preparerForClustersUpdate(ctx context.Context, id ClusterId, input Cluster) (*http.Request, error) {
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

// senderForClustersUpdate sends the ClustersUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c EventHubsClustersClient) senderForClustersUpdate(ctx context.Context, req *http.Request) (future ClustersUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
