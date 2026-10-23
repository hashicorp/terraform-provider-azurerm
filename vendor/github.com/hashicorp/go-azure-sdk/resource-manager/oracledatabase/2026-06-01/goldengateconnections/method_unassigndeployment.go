package goldengateconnections

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

type UnassignDeploymentOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *AssignedDeployment
}

// UnassignDeployment ...
func (c GoldenGateConnectionsClient) UnassignDeployment(ctx context.Context, id GoldenGateConnectionId, input AssignUnassignDeployment) (result UnassignDeploymentOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/unassignDeployment", id.ID()),
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

// UnassignDeploymentThenPoll performs UnassignDeployment then polls until it's completed
func (c GoldenGateConnectionsClient) UnassignDeploymentThenPoll(ctx context.Context, id GoldenGateConnectionId, input AssignUnassignDeployment) error {
	return c.UnassignDeploymentCallbackThenPoll(ctx, id, input, nil)
}

// UnassignDeploymentCallbackThenPoll performs UnassignDeployment, runs the optional callback function, then polls until it's completed
func (c GoldenGateConnectionsClient) UnassignDeploymentCallbackThenPoll(ctx context.Context, id GoldenGateConnectionId, input AssignUnassignDeployment, callback func() error) error {
	result, err := c.UnassignDeployment(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing UnassignDeployment: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after UnassignDeployment: %+v", err)
	}

	return nil
}
