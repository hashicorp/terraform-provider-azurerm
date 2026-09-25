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

type SqlResourcesListSqlTriggersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SqlTriggerGetResults
}

type SqlResourcesListSqlTriggersCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SqlTriggerGetResults
}

type SqlResourcesListSqlTriggersCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *SqlResourcesListSqlTriggersCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// SqlResourcesListSqlTriggers ...
func (c OpenapisClient) SqlResourcesListSqlTriggers(ctx context.Context, id ContainerId) (result SqlResourcesListSqlTriggersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &SqlResourcesListSqlTriggersCustomPager{},
		Path:       fmt.Sprintf("%s/triggers", id.ID()),
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
		Values *[]SqlTriggerGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SqlResourcesListSqlTriggersComplete retrieves all the results into a single object
func (c OpenapisClient) SqlResourcesListSqlTriggersComplete(ctx context.Context, id ContainerId) (SqlResourcesListSqlTriggersCompleteResult, error) {
	return c.SqlResourcesListSqlTriggersCompleteMatchingPredicate(ctx, id, SqlTriggerGetResultsOperationPredicate{})
}

// SqlResourcesListSqlTriggersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) SqlResourcesListSqlTriggersCompleteMatchingPredicate(ctx context.Context, id ContainerId, predicate SqlTriggerGetResultsOperationPredicate) (result SqlResourcesListSqlTriggersCompleteResult, err error) {
	items := make([]SqlTriggerGetResults, 0)

	resp, err := c.SqlResourcesListSqlTriggers(ctx, id)
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

	result = SqlResourcesListSqlTriggersCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
