package workspaces

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

type FailoverOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
}

// Failover ...
func (c WorkspacesClient) Failover(ctx context.Context, id LocationWorkspaceId) (result FailoverOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/failover", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
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

// FailoverThenPoll performs Failover then polls until it's completed
func (c WorkspacesClient) FailoverThenPoll(ctx context.Context, id LocationWorkspaceId) error {
	return c.FailoverCallbackThenPoll(ctx, id, nil)
}

// FailoverCallbackThenPoll performs Failover, runs the optional callback function, then polls until it's completed
func (c WorkspacesClient) FailoverCallbackThenPoll(ctx context.Context, id LocationWorkspaceId, callback func() error) error {
	result, err := c.Failover(ctx, id)
	if err != nil {
		return fmt.Errorf("performing Failover: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after Failover: %+v", err)
	}

	return nil
}
