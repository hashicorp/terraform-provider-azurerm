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

type SqlResourcesListSqlUserDefinedFunctionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SqlUserDefinedFunctionGetResults
}

type SqlResourcesListSqlUserDefinedFunctionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SqlUserDefinedFunctionGetResults
}

type SqlResourcesListSqlUserDefinedFunctionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SqlResourcesListSqlUserDefinedFunctionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SqlResourcesListSqlUserDefinedFunctions ...
func (c OpenapisClient) SqlResourcesListSqlUserDefinedFunctions(ctx context.Context, id ContainerId) (result SqlResourcesListSqlUserDefinedFunctionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SqlResourcesListSqlUserDefinedFunctionsCustomPager{},
		Path:       fmt.Sprintf("%s/userDefinedFunctions", id.ID()),
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
		Values *[]SqlUserDefinedFunctionGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SqlResourcesListSqlUserDefinedFunctionsComplete retrieves all the results into a single object
func (c OpenapisClient) SqlResourcesListSqlUserDefinedFunctionsComplete(ctx context.Context, id ContainerId) (SqlResourcesListSqlUserDefinedFunctionsCompleteResult, error) {
	return c.SqlResourcesListSqlUserDefinedFunctionsCompleteMatchingPredicate(ctx, id, SqlUserDefinedFunctionGetResultsOperationPredicate{})
}

// SqlResourcesListSqlUserDefinedFunctionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) SqlResourcesListSqlUserDefinedFunctionsCompleteMatchingPredicate(ctx context.Context, id ContainerId, predicate SqlUserDefinedFunctionGetResultsOperationPredicate) (result SqlResourcesListSqlUserDefinedFunctionsCompleteResult, err error) {
	items := make([]SqlUserDefinedFunctionGetResults, 0)

	resp, err := c.SqlResourcesListSqlUserDefinedFunctions(ctx, id)
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

	result = SqlResourcesListSqlUserDefinedFunctionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
