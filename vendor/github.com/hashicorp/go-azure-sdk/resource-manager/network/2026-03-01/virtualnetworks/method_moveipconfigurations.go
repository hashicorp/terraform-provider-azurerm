package virtualnetworks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MoveIPConfigurationsOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// MoveIPConfigurations ...
func (c VirtualNetworksClient) MoveIPConfigurations(ctx context.Context, id commonids.VirtualNetworkId, input MoveIPConfigurationsRequest) (result MoveIPConfigurationsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/moveIpConfigurations", id.ID()),
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

// MoveIPConfigurationsThenPoll performs MoveIPConfigurations then polls until it's completed
func (c VirtualNetworksClient) MoveIPConfigurationsThenPoll(ctx context.Context, id commonids.VirtualNetworkId, input MoveIPConfigurationsRequest) error {
	return c.MoveIPConfigurationsCallbackThenPoll(ctx, id, input, nil)
}

// MoveIPConfigurationsCallbackThenPoll performs MoveIPConfigurations, runs the optional callback function, then polls until it's completed
func (c VirtualNetworksClient) MoveIPConfigurationsCallbackThenPoll(ctx context.Context, id commonids.VirtualNetworkId, input MoveIPConfigurationsRequest, callback func() error) error {
	result, err := c.MoveIPConfigurations(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MoveIPConfigurations: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after MoveIPConfigurations: %+v", err)
	}

	return nil
}
