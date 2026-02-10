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

type StartPacketCaptureOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *string
}

// StartPacketCapture ...
func (c VirtualNetworkGatewaysClient) StartPacketCapture(ctx context.Context, id VirtualNetworkGatewayId, input VpnPacketCaptureStartParameters) (result StartPacketCaptureOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/startPacketCapture", id.ID()),
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

// StartPacketCaptureThenPoll performs StartPacketCapture then polls until it's completed
func (c VirtualNetworkGatewaysClient) StartPacketCaptureThenPoll(ctx context.Context, id VirtualNetworkGatewayId, input VpnPacketCaptureStartParameters) error {
	result, err := c.StartPacketCapture(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing StartPacketCapture: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after StartPacketCapture: %+v", err)
	}

	return nil
}
