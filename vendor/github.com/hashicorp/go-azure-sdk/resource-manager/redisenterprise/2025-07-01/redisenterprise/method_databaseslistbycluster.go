package redisenterprise

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabasesListByClusterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Database
}

type DatabasesListByClusterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Database
}

type DatabasesListByClusterCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DatabasesListByClusterCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DatabasesListByCluster ...
func (c RedisEnterpriseClient) DatabasesListByCluster(ctx context.Context, id RedisEnterpriseId) (result DatabasesListByClusterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DatabasesListByClusterCustomPager{},
		Path:       fmt.Sprintf("%s/databases", id.ID()),
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
		Values *[]Database `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DatabasesListByClusterComplete retrieves all the results into a single object
func (c RedisEnterpriseClient) DatabasesListByClusterComplete(ctx context.Context, id RedisEnterpriseId) (DatabasesListByClusterCompleteResult, error) {
	return c.DatabasesListByClusterCompleteMatchingPredicate(ctx, id, DatabaseOperationPredicate{})
}

// DatabasesListByClusterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RedisEnterpriseClient) DatabasesListByClusterCompleteMatchingPredicate(ctx context.Context, id RedisEnterpriseId, predicate DatabaseOperationPredicate) (result DatabasesListByClusterCompleteResult, err error) {
	items := make([]Database, 0)

	resp, err := c.DatabasesListByCluster(ctx, id)
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

	result = DatabasesListByClusterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
