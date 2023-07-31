package topicsauthorizationrule

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicsListAuthorizationRulesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SBAuthorizationRule
}

type TopicsListAuthorizationRulesCompleteResult struct {
	Items []SBAuthorizationRule
}

// TopicsListAuthorizationRules ...
func (c TopicsAuthorizationRuleClient) TopicsListAuthorizationRules(ctx context.Context, id TopicId) (result TopicsListAuthorizationRulesOperationResponse, err error) {
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

// TopicsListAuthorizationRulesComplete retrieves all the results into a single object
func (c TopicsAuthorizationRuleClient) TopicsListAuthorizationRulesComplete(ctx context.Context, id TopicId) (TopicsListAuthorizationRulesCompleteResult, error) {
	return c.TopicsListAuthorizationRulesCompleteMatchingPredicate(ctx, id, SBAuthorizationRuleOperationPredicate{})
}

// TopicsListAuthorizationRulesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c TopicsAuthorizationRuleClient) TopicsListAuthorizationRulesCompleteMatchingPredicate(ctx context.Context, id TopicId, predicate SBAuthorizationRuleOperationPredicate) (result TopicsListAuthorizationRulesCompleteResult, err error) {
	items := make([]SBAuthorizationRule, 0)

	resp, err := c.TopicsListAuthorizationRules(ctx, id)
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

	result = TopicsListAuthorizationRulesCompleteResult{
		Items: items,
	}
	return
}
