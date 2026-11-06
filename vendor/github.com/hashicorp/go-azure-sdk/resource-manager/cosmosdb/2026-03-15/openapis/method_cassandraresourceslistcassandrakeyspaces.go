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

type CassandraResourcesListCassandraKeyspacesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CassandraKeyspaceGetResults
}

type CassandraResourcesListCassandraKeyspacesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CassandraKeyspaceGetResults
}

type CassandraResourcesListCassandraKeyspacesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CassandraResourcesListCassandraKeyspacesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CassandraResourcesListCassandraKeyspaces ...
func (c OpenapisClient) CassandraResourcesListCassandraKeyspaces(ctx context.Context, id DatabaseAccountId) (result CassandraResourcesListCassandraKeyspacesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CassandraResourcesListCassandraKeyspacesCustomPager{},
		Path:       fmt.Sprintf("%s/cassandraKeyspaces", id.ID()),
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
		Values *[]CassandraKeyspaceGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CassandraResourcesListCassandraKeyspacesComplete retrieves all the results into a single object
func (c OpenapisClient) CassandraResourcesListCassandraKeyspacesComplete(ctx context.Context, id DatabaseAccountId) (CassandraResourcesListCassandraKeyspacesCompleteResult, error) {
	return c.CassandraResourcesListCassandraKeyspacesCompleteMatchingPredicate(ctx, id, CassandraKeyspaceGetResultsOperationPredicate{})
}

// CassandraResourcesListCassandraKeyspacesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) CassandraResourcesListCassandraKeyspacesCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate CassandraKeyspaceGetResultsOperationPredicate) (result CassandraResourcesListCassandraKeyspacesCompleteResult, err error) {
	items := make([]CassandraKeyspaceGetResults, 0)

	resp, err := c.CassandraResourcesListCassandraKeyspaces(ctx, id)
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

	result = CassandraResourcesListCassandraKeyspacesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
