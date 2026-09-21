package ipfirewallrules

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

type ReplaceAllOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ReplaceAllFirewallRulesOperationResponse
}

// ReplaceAll ...
func (c IPFirewallRulesClient) ReplaceAll(ctx context.Context, id WorkspaceId, input ReplaceAllIPFirewallRulesRequest) (result ReplaceAllOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/replaceAllIpFirewallRules", id.ID()),
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

// ReplaceAllThenPoll performs ReplaceAll then polls until it's completed
func (c IPFirewallRulesClient) ReplaceAllThenPoll(ctx context.Context, id WorkspaceId, input ReplaceAllIPFirewallRulesRequest) error {
	return c.ReplaceAllCallbackThenPoll(ctx, id, input, nil)
}

// ReplaceAllCallbackThenPoll performs ReplaceAll, runs the optional callback function, then polls until it's completed
func (c IPFirewallRulesClient) ReplaceAllCallbackThenPoll(ctx context.Context, id WorkspaceId, input ReplaceAllIPFirewallRulesRequest, callback func() error) error {
	result, err := c.ReplaceAll(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ReplaceAll: %+v", err)
	}

	if callback != nil {
		if err := callback(); err != nil {
			return fmt.Errorf("executing callback function: %+v", err)
		}
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after ReplaceAll: %+v", err)
	}

	return nil
}
