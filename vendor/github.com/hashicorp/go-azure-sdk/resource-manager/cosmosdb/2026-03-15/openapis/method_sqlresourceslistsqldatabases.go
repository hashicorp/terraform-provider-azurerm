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

type SqlResourcesListSqlDatabasesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SqlDatabaseGetResults
}

type SqlResourcesListSqlDatabasesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SqlDatabaseGetResults
}

type SqlResourcesListSqlDatabasesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SqlResourcesListSqlDatabasesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SqlResourcesListSqlDatabases ...
func (c OpenapisClient) SqlResourcesListSqlDatabases(ctx context.Context, id DatabaseAccountId) (result SqlResourcesListSqlDatabasesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SqlResourcesListSqlDatabasesCustomPager{},
		Path:       fmt.Sprintf("%s/sqlDatabases", id.ID()),
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
		Values *[]SqlDatabaseGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SqlResourcesListSqlDatabasesComplete retrieves all the results into a single object
func (c OpenapisClient) SqlResourcesListSqlDatabasesComplete(ctx context.Context, id DatabaseAccountId) (SqlResourcesListSqlDatabasesCompleteResult, error) {
	return c.SqlResourcesListSqlDatabasesCompleteMatchingPredicate(ctx, id, SqlDatabaseGetResultsOperationPredicate{})
}

// SqlResourcesListSqlDatabasesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) SqlResourcesListSqlDatabasesCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate SqlDatabaseGetResultsOperationPredicate) (result SqlResourcesListSqlDatabasesCompleteResult, err error) {
	items := make([]SqlDatabaseGetResults, 0)

	resp, err := c.SqlResourcesListSqlDatabases(ctx, id)
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

	result = SqlResourcesListSqlDatabasesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
