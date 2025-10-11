package azurefirewalls

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

type PacketCaptureOperationOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *AzureFirewallPacketCaptureResponse
}

// PacketCaptureOperation ...
func (c AzureFirewallsClient) PacketCaptureOperation(ctx context.Context, id AzureFirewallId, input FirewallPacketCaptureParameters) (result PacketCaptureOperationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/packetCaptureOperation", id.ID()),
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

// PacketCaptureOperationThenPoll performs PacketCaptureOperation then polls until it's completed
func (c AzureFirewallsClient) PacketCaptureOperationThenPoll(ctx context.Context, id AzureFirewallId, input FirewallPacketCaptureParameters) error {
	result, err := c.PacketCaptureOperation(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing PacketCaptureOperation: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after PacketCaptureOperation: %+v", err)
	}

	return nil
}
