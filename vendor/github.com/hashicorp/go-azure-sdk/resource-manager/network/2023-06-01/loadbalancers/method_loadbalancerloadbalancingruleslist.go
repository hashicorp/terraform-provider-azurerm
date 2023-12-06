package loadbalancers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancerLoadBalancingRulesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LoadBalancingRule
}

type LoadBalancerLoadBalancingRulesListCompleteResult struct {
	Items []LoadBalancingRule
}

// LoadBalancerLoadBalancingRulesList ...
func (c LoadBalancersClient) LoadBalancerLoadBalancingRulesList(ctx context.Context, id ProviderLoadBalancerId) (result LoadBalancerLoadBalancingRulesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/loadBalancingRules", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]LoadBalancingRule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LoadBalancerLoadBalancingRulesListComplete retrieves all the results into a single object
func (c LoadBalancersClient) LoadBalancerLoadBalancingRulesListComplete(ctx context.Context, id ProviderLoadBalancerId) (LoadBalancerLoadBalancingRulesListCompleteResult, error) {
	return c.LoadBalancerLoadBalancingRulesListCompleteMatchingPredicate(ctx, id, LoadBalancingRuleOperationPredicate{})
}

// LoadBalancerLoadBalancingRulesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LoadBalancersClient) LoadBalancerLoadBalancingRulesListCompleteMatchingPredicate(ctx context.Context, id ProviderLoadBalancerId, predicate LoadBalancingRuleOperationPredicate) (result LoadBalancerLoadBalancingRulesListCompleteResult, err error) {
	items := make([]LoadBalancingRule, 0)

	resp, err := c.LoadBalancerLoadBalancingRulesList(ctx, id)
	if err != nil {
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = LoadBalancerLoadBalancingRulesListCompleteResult{
		Items: items,
	}
	return
}
