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

type AssignConnectionOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *AssignedConnection
}

// AssignConnection ...
func (c GoldenGateDeploymentsClient) AssignConnection(ctx context.Context, id GoldenGateDeploymentId, input AssignUnassignConnection) (result AssignConnectionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/assignConnection", id.ID()),
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

// AssignConnectionThenPoll performs AssignConnection then polls until it's completed
func (c GoldenGateDeploymentsClient) AssignConnectionThenPoll(ctx context.Context, id GoldenGateDeploymentId, input AssignUnassignConnection) error {
	return c.AssignConnectionCallbackThenPoll(ctx, id, input, nil)
}

// AssignConnectionCallbackThenPoll performs AssignConnection, runs the optional callback function, then polls until it's completed
func (c GoldenGateDeploymentsClient) AssignConnectionCallbackThenPoll(ctx context.Context, id GoldenGateDeploymentId, input AssignUnassignConnection, callback func() error) error {
	result, err := c.AssignConnection(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing AssignConnection: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after AssignConnection: %+v", err)
	}

	return nil
}
