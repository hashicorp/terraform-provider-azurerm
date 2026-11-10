package rules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByRuleSetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Rule
}

type ListByRuleSetCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Rule
}

type ListByRuleSetCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByRuleSetCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByRuleSet ...
func (c RulesClient) ListByRuleSet(ctx context.Context, id RuleSetId) (result ListByRuleSetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByRuleSetCustomPager{},
		Path:       fmt.Sprintf("%s/rules", id.ID()),
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
		Values *[]Rule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByRuleSetComplete retrieves all the results into a single object
func (c RulesClient) ListByRuleSetComplete(ctx context.Context, id RuleSetId) (ListByRuleSetCompleteResult, error) {
	return c.ListByRuleSetCompleteMatchingPredicate(ctx, id, RuleOperationPredicate{})
}

// ListByRuleSetCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RulesClient) ListByRuleSetCompleteMatchingPredicate(ctx context.Context, id RuleSetId, predicate RuleOperationPredicate) (result ListByRuleSetCompleteResult, err error) {
	items := make([]Rule, 0)

	resp, err := c.ListByRuleSet(ctx, id)
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

	result = ListByRuleSetCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
