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

type SqlResourcesListSqlRoleDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SqlRoleDefinitionGetResults
}

type SqlResourcesListSqlRoleDefinitionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SqlRoleDefinitionGetResults
}

type SqlResourcesListSqlRoleDefinitionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SqlResourcesListSqlRoleDefinitionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SqlResourcesListSqlRoleDefinitions ...
func (c OpenapisClient) SqlResourcesListSqlRoleDefinitions(ctx context.Context, id DatabaseAccountId) (result SqlResourcesListSqlRoleDefinitionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SqlResourcesListSqlRoleDefinitionsCustomPager{},
		Path:       fmt.Sprintf("%s/sqlRoleDefinitions", id.ID()),
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
		Values *[]SqlRoleDefinitionGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SqlResourcesListSqlRoleDefinitionsComplete retrieves all the results into a single object
func (c OpenapisClient) SqlResourcesListSqlRoleDefinitionsComplete(ctx context.Context, id DatabaseAccountId) (SqlResourcesListSqlRoleDefinitionsCompleteResult, error) {
	return c.SqlResourcesListSqlRoleDefinitionsCompleteMatchingPredicate(ctx, id, SqlRoleDefinitionGetResultsOperationPredicate{})
}

// SqlResourcesListSqlRoleDefinitionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) SqlResourcesListSqlRoleDefinitionsCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate SqlRoleDefinitionGetResultsOperationPredicate) (result SqlResourcesListSqlRoleDefinitionsCompleteResult, err error) {
	items := make([]SqlRoleDefinitionGetResults, 0)

	resp, err := c.SqlResourcesListSqlRoleDefinitions(ctx, id)
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

	result = SqlResourcesListSqlRoleDefinitionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
