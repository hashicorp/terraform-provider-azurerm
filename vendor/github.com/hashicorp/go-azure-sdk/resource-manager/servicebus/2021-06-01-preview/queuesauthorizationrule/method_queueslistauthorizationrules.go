package queuesauthorizationrule

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueuesListAuthorizationRulesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SBAuthorizationRule
}

type QueuesListAuthorizationRulesCompleteResult struct {
	Items []SBAuthorizationRule
}

// QueuesListAuthorizationRules ...
func (c QueuesAuthorizationRuleClient) QueuesListAuthorizationRules(ctx context.Context, id QueueId) (result QueuesListAuthorizationRulesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/authorizationRules", id.ID()),
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
		Values *[]SBAuthorizationRule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// QueuesListAuthorizationRulesComplete retrieves all the results into a single object
func (c QueuesAuthorizationRuleClient) QueuesListAuthorizationRulesComplete(ctx context.Context, id QueueId) (QueuesListAuthorizationRulesCompleteResult, error) {
	return c.QueuesListAuthorizationRulesCompleteMatchingPredicate(ctx, id, SBAuthorizationRuleOperationPredicate{})
}

// QueuesListAuthorizationRulesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c QueuesAuthorizationRuleClient) QueuesListAuthorizationRulesCompleteMatchingPredicate(ctx context.Context, id QueueId, predicate SBAuthorizationRuleOperationPredicate) (result QueuesListAuthorizationRulesCompleteResult, err error) {
	items := make([]SBAuthorizationRule, 0)

	resp, err := c.QueuesListAuthorizationRules(ctx, id)
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

	result = QueuesListAuthorizationRulesCompleteResult{
		Items: items,
	}
	return
}
