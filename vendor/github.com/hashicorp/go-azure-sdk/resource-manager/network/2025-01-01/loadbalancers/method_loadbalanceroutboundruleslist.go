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

type LoadBalancerOutboundRulesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]OutboundRule
}

type LoadBalancerOutboundRulesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []OutboundRule
}

type LoadBalancerOutboundRulesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LoadBalancerOutboundRulesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LoadBalancerOutboundRulesList ...
func (c LoadBalancersClient) LoadBalancerOutboundRulesList(ctx context.Context, id ProviderLoadBalancerId) (result LoadBalancerOutboundRulesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &LoadBalancerOutboundRulesListCustomPager{},
		Path:       fmt.Sprintf("%s/outboundRules", id.ID()),
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
		Values *[]OutboundRule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LoadBalancerOutboundRulesListComplete retrieves all the results into a single object
func (c LoadBalancersClient) LoadBalancerOutboundRulesListComplete(ctx context.Context, id ProviderLoadBalancerId) (LoadBalancerOutboundRulesListCompleteResult, error) {
	return c.LoadBalancerOutboundRulesListCompleteMatchingPredicate(ctx, id, OutboundRuleOperationPredicate{})
}

// LoadBalancerOutboundRulesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LoadBalancersClient) LoadBalancerOutboundRulesListCompleteMatchingPredicate(ctx context.Context, id ProviderLoadBalancerId, predicate OutboundRuleOperationPredicate) (result LoadBalancerOutboundRulesListCompleteResult, err error) {
	items := make([]OutboundRule, 0)

	resp, err := c.LoadBalancerOutboundRulesList(ctx, id)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
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

	result = LoadBalancerOutboundRulesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
