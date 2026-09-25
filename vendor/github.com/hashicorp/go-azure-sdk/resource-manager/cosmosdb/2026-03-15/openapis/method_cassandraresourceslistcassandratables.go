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

type CassandraResourcesListCassandraTablesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CassandraTableGetResults
}

type CassandraResourcesListCassandraTablesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CassandraTableGetResults
}

type CassandraResourcesListCassandraTablesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CassandraResourcesListCassandraTablesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CassandraResourcesListCassandraTables ...
func (c OpenapisClient) CassandraResourcesListCassandraTables(ctx context.Context, id CassandraKeyspaceId) (result CassandraResourcesListCassandraTablesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CassandraResourcesListCassandraTablesCustomPager{},
		Path:       fmt.Sprintf("%s/tables", id.ID()),
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
		Values *[]CassandraTableGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CassandraResourcesListCassandraTablesComplete retrieves all the results into a single object
func (c OpenapisClient) CassandraResourcesListCassandraTablesComplete(ctx context.Context, id CassandraKeyspaceId) (CassandraResourcesListCassandraTablesCompleteResult, error) {
	return c.CassandraResourcesListCassandraTablesCompleteMatchingPredicate(ctx, id, CassandraTableGetResultsOperationPredicate{})
}

// CassandraResourcesListCassandraTablesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) CassandraResourcesListCassandraTablesCompleteMatchingPredicate(ctx context.Context, id CassandraKeyspaceId, predicate CassandraTableGetResultsOperationPredicate) (result CassandraResourcesListCassandraTablesCompleteResult, err error) {
	items := make([]CassandraTableGetResults, 0)

	resp, err := c.CassandraResourcesListCassandraTables(ctx, id)
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

	result = CassandraResourcesListCassandraTablesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
