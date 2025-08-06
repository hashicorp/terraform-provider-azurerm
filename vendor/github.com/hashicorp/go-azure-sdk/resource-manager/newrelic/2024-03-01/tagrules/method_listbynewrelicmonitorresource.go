package tagrules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByNewRelicMonitorResourceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TagRule
}

type ListByNewRelicMonitorResourceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TagRule
}

type ListByNewRelicMonitorResourceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByNewRelicMonitorResourceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByNewRelicMonitorResource ...
func (c TagRulesClient) ListByNewRelicMonitorResource(ctx context.Context, id MonitorId) (result ListByNewRelicMonitorResourceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByNewRelicMonitorResourceCustomPager{},
		Path:       fmt.Sprintf("%s/tagRules", id.ID()),
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
		Values *[]TagRule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByNewRelicMonitorResourceComplete retrieves all the results into a single object
func (c TagRulesClient) ListByNewRelicMonitorResourceComplete(ctx context.Context, id MonitorId) (ListByNewRelicMonitorResourceCompleteResult, error) {
	return c.ListByNewRelicMonitorResourceCompleteMatchingPredicate(ctx, id, TagRuleOperationPredicate{})
}

// ListByNewRelicMonitorResourceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c TagRulesClient) ListByNewRelicMonitorResourceCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate TagRuleOperationPredicate) (result ListByNewRelicMonitorResourceCompleteResult, err error) {
	items := make([]TagRule, 0)

	resp, err := c.ListByNewRelicMonitorResource(ctx, id)
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

	result = ListByNewRelicMonitorResourceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
