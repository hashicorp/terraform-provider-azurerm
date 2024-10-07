package routefilterrules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByRouteFilterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RouteFilterRule
}

type ListByRouteFilterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RouteFilterRule
}

type ListByRouteFilterCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByRouteFilterCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByRouteFilter ...
func (c RouteFilterRulesClient) ListByRouteFilter(ctx context.Context, id RouteFilterId) (result ListByRouteFilterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByRouteFilterCustomPager{},
		Path:       fmt.Sprintf("%s/routeFilterRules", id.ID()),
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
		Values *[]RouteFilterRule `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByRouteFilterComplete retrieves all the results into a single object
func (c RouteFilterRulesClient) ListByRouteFilterComplete(ctx context.Context, id RouteFilterId) (ListByRouteFilterCompleteResult, error) {
	return c.ListByRouteFilterCompleteMatchingPredicate(ctx, id, RouteFilterRuleOperationPredicate{})
}

// ListByRouteFilterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RouteFilterRulesClient) ListByRouteFilterCompleteMatchingPredicate(ctx context.Context, id RouteFilterId, predicate RouteFilterRuleOperationPredicate) (result ListByRouteFilterCompleteResult, err error) {
	items := make([]RouteFilterRule, 0)

	resp, err := c.ListByRouteFilter(ctx, id)
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

	result = ListByRouteFilterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
