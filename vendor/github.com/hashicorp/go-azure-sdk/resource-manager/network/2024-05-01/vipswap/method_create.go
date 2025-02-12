package vipswap

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

type CreateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// Create ...
func (c VipSwapClient) Create(ctx context.Context, id CloudServiceId, input SwapResource) (result CreateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Network/cloudServiceSlots/swap", id.ID()),
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

// CreateThenPoll performs Create then polls until it's completed
func (c VipSwapClient) CreateThenPoll(ctx context.Context, id CloudServiceId, input SwapResource) error {
	result, err := c.Create(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Create: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Create: %+v", err)
	}

	return nil
}
