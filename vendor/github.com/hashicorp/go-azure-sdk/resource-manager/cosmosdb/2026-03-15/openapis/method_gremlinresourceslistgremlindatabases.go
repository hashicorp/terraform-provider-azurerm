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

type GremlinResourcesListGremlinDatabasesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GremlinDatabaseGetResults
}

type GremlinResourcesListGremlinDatabasesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GremlinDatabaseGetResults
}

type GremlinResourcesListGremlinDatabasesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GremlinResourcesListGremlinDatabasesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GremlinResourcesListGremlinDatabases ...
func (c OpenapisClient) GremlinResourcesListGremlinDatabases(ctx context.Context, id DatabaseAccountId) (result GremlinResourcesListGremlinDatabasesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GremlinResourcesListGremlinDatabasesCustomPager{},
		Path:       fmt.Sprintf("%s/gremlinDatabases", id.ID()),
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
		Values *[]GremlinDatabaseGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GremlinResourcesListGremlinDatabasesComplete retrieves all the results into a single object
func (c OpenapisClient) GremlinResourcesListGremlinDatabasesComplete(ctx context.Context, id DatabaseAccountId) (GremlinResourcesListGremlinDatabasesCompleteResult, error) {
	return c.GremlinResourcesListGremlinDatabasesCompleteMatchingPredicate(ctx, id, GremlinDatabaseGetResultsOperationPredicate{})
}

// GremlinResourcesListGremlinDatabasesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) GremlinResourcesListGremlinDatabasesCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate GremlinDatabaseGetResultsOperationPredicate) (result GremlinResourcesListGremlinDatabasesCompleteResult, err error) {
	items := make([]GremlinDatabaseGetResults, 0)

	resp, err := c.GremlinResourcesListGremlinDatabases(ctx, id)
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

	result = GremlinResourcesListGremlinDatabasesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
