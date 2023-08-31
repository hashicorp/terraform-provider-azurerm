package redis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallRulesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RedisFirewallRule
}

type FirewallRulesListCompleteResult struct {
	Items []RedisFirewallRule
}

// FirewallRulesList ...
func (c RedisClient) FirewallRulesList(ctx context.Context, id RediId) (result FirewallRulesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/firewallRules", id.ID()),
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
		Values *[]RedisFirewallRule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// FirewallRulesListComplete retrieves all the results into a single object
func (c RedisClient) FirewallRulesListComplete(ctx context.Context, id RediId) (FirewallRulesListCompleteResult, error) {
	return c.FirewallRulesListCompleteMatchingPredicate(ctx, id, RedisFirewallRuleOperationPredicate{})
}

// FirewallRulesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RedisClient) FirewallRulesListCompleteMatchingPredicate(ctx context.Context, id RediId, predicate RedisFirewallRuleOperationPredicate) (result FirewallRulesListCompleteResult, err error) {
	items := make([]RedisFirewallRule, 0)

	resp, err := c.FirewallRulesList(ctx, id)
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

	result = FirewallRulesListCompleteResult{
		Items: items,
	}
	return
}
