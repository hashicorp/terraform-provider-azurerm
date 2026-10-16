package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GremlinResourcesListGremlinGraphsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GremlinGraphGetResults
}

type GremlinResourcesListGremlinGraphsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GremlinGraphGetResults
}

type GremlinResourcesListGremlinGraphsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GremlinResourcesListGremlinGraphsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GremlinResourcesListGremlinGraphs ...
func (c OpenapisClient) GremlinResourcesListGremlinGraphs(ctx context.Context, id GremlinDatabaseId) (result GremlinResourcesListGremlinGraphsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GremlinResourcesListGremlinGraphsCustomPager{},
		Path:       fmt.Sprintf("%s/graphs", id.ID()),
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
		Values *[]GremlinGraphGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GremlinResourcesListGremlinGraphsComplete retrieves all the results into a single object
func (c OpenapisClient) GremlinResourcesListGremlinGraphsComplete(ctx context.Context, id GremlinDatabaseId) (GremlinResourcesListGremlinGraphsCompleteResult, error) {
	return c.GremlinResourcesListGremlinGraphsCompleteMatchingPredicate(ctx, id, GremlinGraphGetResultsOperationPredicate{})
}

// GremlinResourcesListGremlinGraphsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) GremlinResourcesListGremlinGraphsCompleteMatchingPredicate(ctx context.Context, id GremlinDatabaseId, predicate GremlinGraphGetResultsOperationPredicate) (result GremlinResourcesListGremlinGraphsCompleteResult, err error) {
	items := make([]GremlinGraphGetResults, 0)

	resp, err := c.GremlinResourcesListGremlinGraphs(ctx, id)
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

	result = GremlinResourcesListGremlinGraphsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
