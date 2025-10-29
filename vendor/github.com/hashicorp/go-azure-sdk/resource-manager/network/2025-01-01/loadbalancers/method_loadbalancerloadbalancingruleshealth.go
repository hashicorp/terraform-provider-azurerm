package loadbalancers

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

type LoadBalancerLoadBalancingRulesHealthOperationResponse struct {
	Poller       pollers.Poller
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *LoadBalancerHealthPerRule
}

// LoadBalancerLoadBalancingRulesHealth ...
func (c LoadBalancersClient) LoadBalancerLoadBalancingRulesHealth(ctx context.Context, id LoadBalancingRuleId) (result LoadBalancerLoadBalancingRulesHealthOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/health", id.ID()),
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

// LoadBalancerLoadBalancingRulesHealthThenPoll performs LoadBalancerLoadBalancingRulesHealth then polls until it's completed
func (c LoadBalancersClient) LoadBalancerLoadBalancingRulesHealthThenPoll(ctx context.Context, id LoadBalancingRuleId) error {
	result, err := c.LoadBalancerLoadBalancingRulesHealth(ctx, id)
	if err != nil {
		return fmt.Errorf("performing LoadBalancerLoadBalancingRulesHealth: %+v", err)
	}

	if err := result.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("polling after LoadBalancerLoadBalancingRulesHealth: %+v", err)
	}

	return nil
}
