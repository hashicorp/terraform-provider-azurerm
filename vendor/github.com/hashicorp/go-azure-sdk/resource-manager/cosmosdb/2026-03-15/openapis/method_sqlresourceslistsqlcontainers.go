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

type SqlResourcesListSqlContainersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SqlContainerGetResults
}

type SqlResourcesListSqlContainersCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SqlContainerGetResults
}

type SqlResourcesListSqlContainersCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SqlResourcesListSqlContainersCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SqlResourcesListSqlContainers ...
func (c OpenapisClient) SqlResourcesListSqlContainers(ctx context.Context, id SqlDatabaseId) (result SqlResourcesListSqlContainersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SqlResourcesListSqlContainersCustomPager{},
		Path:       fmt.Sprintf("%s/containers", id.ID()),
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
		Values *[]SqlContainerGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SqlResourcesListSqlContainersComplete retrieves all the results into a single object
func (c OpenapisClient) SqlResourcesListSqlContainersComplete(ctx context.Context, id SqlDatabaseId) (SqlResourcesListSqlContainersCompleteResult, error) {
	return c.SqlResourcesListSqlContainersCompleteMatchingPredicate(ctx, id, SqlContainerGetResultsOperationPredicate{})
}

// SqlResourcesListSqlContainersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) SqlResourcesListSqlContainersCompleteMatchingPredicate(ctx context.Context, id SqlDatabaseId, predicate SqlContainerGetResultsOperationPredicate) (result SqlResourcesListSqlContainersCompleteResult, err error) {
	items := make([]SqlContainerGetResults, 0)

	resp, err := c.SqlResourcesListSqlContainers(ctx, id)
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

	result = SqlResourcesListSqlContainersCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
