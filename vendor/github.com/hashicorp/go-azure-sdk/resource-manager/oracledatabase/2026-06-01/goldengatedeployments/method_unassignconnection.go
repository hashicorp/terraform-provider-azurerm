package goldengatedeployments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UnassignConnectionOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *AssignedConnection
}

// UnassignConnection ...
func (c GoldenGateDeploymentsClient) UnassignConnection(ctx context.Context, id GoldenGateDeploymentId, input AssignUnassignConnection) (result UnassignConnectionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/unassignConnection", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	result.Poller, err = resourcemanager.PollerFromResponse(resp, c.Client)
	if err != nil {
		return
	}

	return
}

// UnassignConnectionThenPoll performs UnassignConnection then polls until it's completed
func (c GoldenGateDeploymentsClient) UnassignConnectionThenPoll(ctx context.Context, id GoldenGateDeploymentId, input AssignUnassignConnection) error {
	return c.UnassignConnectionCallbackThenPoll(ctx, id, input, nil)
}

// UnassignConnectionCallbackThenPoll performs UnassignConnection, runs the optional callback function, then polls until it's completed
func (c GoldenGateDeploymentsClient) UnassignConnectionCallbackThenPoll(ctx context.Context, id GoldenGateDeploymentId, input AssignUnassignConnection, callback func() error) error {
	result, err := c.UnassignConnection(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing UnassignConnection: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after UnassignConnection: %+v", err)
	}

	return nil
}
