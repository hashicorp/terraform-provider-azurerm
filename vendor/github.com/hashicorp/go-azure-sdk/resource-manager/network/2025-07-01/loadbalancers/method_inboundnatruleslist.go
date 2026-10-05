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

type InboundNatRulesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]InboundNatRule
}

type InboundNatRulesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []InboundNatRule
}

type InboundNatRulesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *InboundNatRulesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// InboundNatRulesList ...
func (c LoadBalancersClient) InboundNatRulesList(ctx context.Context, id ProviderLoadBalancerId) (result InboundNatRulesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &InboundNatRulesListCustomPager{},
		Path:       fmt.Sprintf("%s/inboundNatRules", id.ID()),
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
		Values *[]InboundNatRule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// InboundNatRulesListComplete retrieves all the results into a single object
func (c LoadBalancersClient) InboundNatRulesListComplete(ctx context.Context, id ProviderLoadBalancerId) (InboundNatRulesListCompleteResult, error) {
	return c.InboundNatRulesListCompleteMatchingPredicate(ctx, id, InboundNatRuleOperationPredicate{})
}

// InboundNatRulesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LoadBalancersClient) InboundNatRulesListCompleteMatchingPredicate(ctx context.Context, id ProviderLoadBalancerId, predicate InboundNatRuleOperationPredicate) (result InboundNatRulesListCompleteResult, err error) {
	items := make([]InboundNatRule, 0)

	resp, err := c.InboundNatRulesList(ctx, id)
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

	result = InboundNatRulesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
