package deploymentstacksatresourcegroup

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

type DeploymentStacksCreateOrUpdateAtResourceGroupOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *DeploymentStack
}

// DeploymentStacksCreateOrUpdateAtResourceGroup ...
func (c DeploymentStacksAtResourceGroupClient) DeploymentStacksCreateOrUpdateAtResourceGroup(ctx context.Context, id ProviderDeploymentStackId, input DeploymentStack) (result DeploymentStacksCreateOrUpdateAtResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       id.ID(),
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

// DeploymentStacksCreateOrUpdateAtResourceGroupThenPoll performs DeploymentStacksCreateOrUpdateAtResourceGroup then polls until it's completed
func (c DeploymentStacksAtResourceGroupClient) DeploymentStacksCreateOrUpdateAtResourceGroupThenPoll(ctx context.Context, id ProviderDeploymentStackId, input DeploymentStack) error {
	return c.DeploymentStacksCreateOrUpdateAtResourceGroupCallbackThenPoll(ctx, id, input, nil)
}

// DeploymentStacksCreateOrUpdateAtResourceGroupCallbackThenPoll performs DeploymentStacksCreateOrUpdateAtResourceGroup, runs the optional callback function, then polls until it's completed
func (c DeploymentStacksAtResourceGroupClient) DeploymentStacksCreateOrUpdateAtResourceGroupCallbackThenPoll(ctx context.Context, id ProviderDeploymentStackId, input DeploymentStack, callback func() error) error {
	result, err := c.DeploymentStacksCreateOrUpdateAtResourceGroup(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing DeploymentStacksCreateOrUpdateAtResourceGroup: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after DeploymentStacksCreateOrUpdateAtResourceGroup: %+v", err)
	}

	return nil
}
