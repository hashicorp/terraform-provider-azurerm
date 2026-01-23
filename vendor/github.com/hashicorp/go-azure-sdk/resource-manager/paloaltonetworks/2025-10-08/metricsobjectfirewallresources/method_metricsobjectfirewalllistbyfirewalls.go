package metricsobjectfirewallresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetricsObjectFirewallListByFirewallsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MetricsObjectFirewallResource
}

type MetricsObjectFirewallListByFirewallsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []MetricsObjectFirewallResource
}

type MetricsObjectFirewallListByFirewallsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MetricsObjectFirewallListByFirewallsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MetricsObjectFirewallListByFirewalls ...
func (c MetricsObjectFirewallResourcesClient) MetricsObjectFirewallListByFirewalls(ctx context.Context, id FirewallId) (result MetricsObjectFirewallListByFirewallsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &MetricsObjectFirewallListByFirewallsCustomPager{},
		Path:       fmt.Sprintf("%s/metrics", id.ID()),
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
		Values *[]MetricsObjectFirewallResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MetricsObjectFirewallListByFirewallsComplete retrieves all the results into a single object
func (c MetricsObjectFirewallResourcesClient) MetricsObjectFirewallListByFirewallsComplete(ctx context.Context, id FirewallId) (MetricsObjectFirewallListByFirewallsCompleteResult, error) {
	return c.MetricsObjectFirewallListByFirewallsCompleteMatchingPredicate(ctx, id, MetricsObjectFirewallResourceOperationPredicate{})
}

// MetricsObjectFirewallListByFirewallsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c MetricsObjectFirewallResourcesClient) MetricsObjectFirewallListByFirewallsCompleteMatchingPredicate(ctx context.Context, id FirewallId, predicate MetricsObjectFirewallResourceOperationPredicate) (result MetricsObjectFirewallListByFirewallsCompleteResult, err error) {
	items := make([]MetricsObjectFirewallResource, 0)

	resp, err := c.MetricsObjectFirewallListByFirewalls(ctx, id)
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

	result = MetricsObjectFirewallListByFirewallsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
