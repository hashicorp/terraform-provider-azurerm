package virtualnetworkgateways

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

type GeneratevpnclientpackageOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *string
}

// Generatevpnclientpackage ...
func (c VirtualNetworkGatewaysClient) Generatevpnclientpackage(ctx context.Context, id VirtualNetworkGatewayId, input VpnClientParameters) (result GeneratevpnclientpackageOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/generatevpnclientpackage", id.ID()),
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

// GeneratevpnclientpackageThenPoll performs Generatevpnclientpackage then polls until it's completed
func (c VirtualNetworkGatewaysClient) GeneratevpnclientpackageThenPoll(ctx context.Context, id VirtualNetworkGatewayId, input VpnClientParameters) error {
	result, err := c.Generatevpnclientpackage(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Generatevpnclientpackage: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Generatevpnclientpackage: %+v", err)
	}

	return nil
}
