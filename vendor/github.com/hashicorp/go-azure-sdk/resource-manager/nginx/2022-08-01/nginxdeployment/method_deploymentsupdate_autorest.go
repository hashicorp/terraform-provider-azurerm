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

type DeploymentsUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// DeploymentsUpdate ...
func (c NginxDeploymentClient) DeploymentsUpdate(ctx context.Context, id NginxDeploymentId, input NginxDeploymentUpdateParameters) (result DeploymentsUpdateOperationResponse, err error) {
	req, err := c.preparerForDeploymentsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForDeploymentsUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// DeploymentsUpdateThenPoll performs DeploymentsUpdate then polls until it's completed
func (c NginxDeploymentClient) DeploymentsUpdateThenPoll(ctx context.Context, id NginxDeploymentId, input NginxDeploymentUpdateParameters) error {
	result, err := c.DeploymentsUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DeploymentsUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after DeploymentsUpdate: %+v", err)
	}

	return nil
}

// preparerForDeploymentsUpdate prepares the DeploymentsUpdate request.
func (c NginxDeploymentClient) preparerForDeploymentsUpdate(ctx context.Context, id NginxDeploymentId, input NginxDeploymentUpdateParameters) (*http.Request, error) {
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

// senderForDeploymentsUpdate sends the DeploymentsUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c NginxDeploymentClient) senderForDeploymentsUpdate(ctx context.Context, req *http.Request) (future DeploymentsUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
