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

type SqlResourcesListSqlStoredProceduresOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SqlStoredProcedureGetResults
}

type SqlResourcesListSqlStoredProceduresCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SqlStoredProcedureGetResults
}

type SqlResourcesListSqlStoredProceduresCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SqlResourcesListSqlStoredProceduresCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SqlResourcesListSqlStoredProcedures ...
func (c OpenapisClient) SqlResourcesListSqlStoredProcedures(ctx context.Context, id ContainerId) (result SqlResourcesListSqlStoredProceduresOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SqlResourcesListSqlStoredProceduresCustomPager{},
		Path:       fmt.Sprintf("%s/storedProcedures", id.ID()),
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
		Values *[]SqlStoredProcedureGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SqlResourcesListSqlStoredProceduresComplete retrieves all the results into a single object
func (c OpenapisClient) SqlResourcesListSqlStoredProceduresComplete(ctx context.Context, id ContainerId) (SqlResourcesListSqlStoredProceduresCompleteResult, error) {
	return c.SqlResourcesListSqlStoredProceduresCompleteMatchingPredicate(ctx, id, SqlStoredProcedureGetResultsOperationPredicate{})
}

// SqlResourcesListSqlStoredProceduresCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) SqlResourcesListSqlStoredProceduresCompleteMatchingPredicate(ctx context.Context, id ContainerId, predicate SqlStoredProcedureGetResultsOperationPredicate) (result SqlResourcesListSqlStoredProceduresCompleteResult, err error) {
	items := make([]SqlStoredProcedureGetResults, 0)

	resp, err := c.SqlResourcesListSqlStoredProcedures(ctx, id)
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

	result = SqlResourcesListSqlStoredProceduresCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
