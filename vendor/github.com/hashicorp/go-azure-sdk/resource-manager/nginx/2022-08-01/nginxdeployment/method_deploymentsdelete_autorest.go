package nginxdeployment

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

type DeploymentsDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DeploymentsDelete ...
func (c NginxDeploymentClient) DeploymentsDelete(ctx context.Context, id NginxDeploymentId) (result DeploymentsDeleteOperationResponse, err error) {
	req, err := c.preparerForDeploymentsDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDeploymentsDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DeploymentsDeleteThenPoll performs DeploymentsDelete then polls until it's completed
func (c NginxDeploymentClient) DeploymentsDeleteThenPoll(ctx context.Context, id NginxDeploymentId) error {
	result, err := c.DeploymentsDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing DeploymentsDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DeploymentsDelete: %+v", err)
	}

	return nil
}

// preparerForDeploymentsDelete prepares the DeploymentsDelete request.
func (c NginxDeploymentClient) preparerForDeploymentsDelete(ctx context.Context, id NginxDeploymentId) (*http.Request, error) {
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

// senderForDeploymentsDelete sends the DeploymentsDelete request. The method will close the
// http.Response Body if it receives an error.
func (c NginxDeploymentClient) senderForDeploymentsDelete(ctx context.Context, req *http.Request) (future DeploymentsDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
