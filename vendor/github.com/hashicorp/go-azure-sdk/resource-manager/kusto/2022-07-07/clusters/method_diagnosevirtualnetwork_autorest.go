package clusters

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

type DiagnoseVirtualNetworkOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DiagnoseVirtualNetwork ...
func (c ClustersClient) DiagnoseVirtualNetwork(ctx context.Context, id ClusterId) (result DiagnoseVirtualNetworkOperationResponse, err error) {
	req, err := c.preparerForDiagnoseVirtualNetwork(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "DiagnoseVirtualNetwork", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDiagnoseVirtualNetwork(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "DiagnoseVirtualNetwork", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DiagnoseVirtualNetworkThenPoll performs DiagnoseVirtualNetwork then polls until it's completed
func (c ClustersClient) DiagnoseVirtualNetworkThenPoll(ctx context.Context, id ClusterId) error {
	result, err := c.DiagnoseVirtualNetwork(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DiagnoseVirtualNetwork: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DiagnoseVirtualNetwork: %+v", err)
	}

	return nil
}

// preparerForDiagnoseVirtualNetwork prepares the DiagnoseVirtualNetwork request.
func (c ClustersClient) preparerForDiagnoseVirtualNetwork(ctx context.Context, id ClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/diagnoseVirtualNetwork", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForDiagnoseVirtualNetwork sends the DiagnoseVirtualNetwork request. The method will close the
// http.Response Body if it receives an error.
func (c ClustersClient) senderForDiagnoseVirtualNetwork(ctx context.Context, req *http.Request) (future DiagnoseVirtualNetworkOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
