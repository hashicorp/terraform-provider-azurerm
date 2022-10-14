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

type DeploymentsCreateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DeploymentsCreate ...
func (c NginxDeploymentClient) DeploymentsCreate(ctx context.Context, id NginxDeploymentId, input NginxDeployment) (result DeploymentsCreateOperationResponse, err error) {
	req, err := c.preparerForDeploymentsCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsCreate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDeploymentsCreate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsCreate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DeploymentsCreateThenPoll performs DeploymentsCreate then polls until it's completed
func (c NginxDeploymentClient) DeploymentsCreateThenPoll(ctx context.Context, id NginxDeploymentId, input NginxDeployment) error {
	result, err := c.DeploymentsCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DeploymentsCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DeploymentsCreate: %+v", err)
	}

	return nil
}

// preparerForDeploymentsCreate prepares the DeploymentsCreate request.
func (c NginxDeploymentClient) preparerForDeploymentsCreate(ctx context.Context, id NginxDeploymentId, input NginxDeployment) (*http.Request, error) {
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

// senderForDeploymentsCreate sends the DeploymentsCreate request. The method will close the
// http.Response Body if it receives an error.
func (c NginxDeploymentClient) senderForDeploymentsCreate(ctx context.Context, req *http.Request) (future DeploymentsCreateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
