package guestagents

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

type GuestAgentCreateOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *GuestAgent
}

// GuestAgentCreate ...
func (c GuestAgentsClient) GuestAgentCreate(ctx context.Context, id commonids.ScopeId, input GuestAgent) (result GuestAgentCreateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		Path:       fmt.Sprintf("%s/providers/Microsoft.AzureStackHCI/virtualMachineInstances/default/guestAgents/default", id.ID()),
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

// GuestAgentCreateThenPoll performs GuestAgentCreate then polls until it's completed
func (c GuestAgentsClient) GuestAgentCreateThenPoll(ctx context.Context, id commonids.ScopeId, input GuestAgent) error {
	result, err := c.GuestAgentCreate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing GuestAgentCreate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after GuestAgentCreate: %+v", err)
	}

	return nil
}
