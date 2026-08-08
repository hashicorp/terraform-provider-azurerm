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

type RestorableGremlinDatabasesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RestorableGremlinDatabaseGetResult
}

type RestorableGremlinDatabasesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RestorableGremlinDatabaseGetResult
}

type RestorableGremlinDatabasesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RestorableGremlinDatabasesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RestorableGremlinDatabasesList ...
func (c OpenapisClient) RestorableGremlinDatabasesList(ctx context.Context, id RestorableDatabaseAccountId) (result RestorableGremlinDatabasesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &RestorableGremlinDatabasesListCustomPager{},
		Path:       fmt.Sprintf("%s/restorableGremlinDatabases", id.ID()),
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
		Values *[]RestorableGremlinDatabaseGetResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RestorableGremlinDatabasesListComplete retrieves all the results into a single object
func (c OpenapisClient) RestorableGremlinDatabasesListComplete(ctx context.Context, id RestorableDatabaseAccountId) (RestorableGremlinDatabasesListCompleteResult, error) {
	return c.RestorableGremlinDatabasesListCompleteMatchingPredicate(ctx, id, RestorableGremlinDatabaseGetResultOperationPredicate{})
}

// RestorableGremlinDatabasesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) RestorableGremlinDatabasesListCompleteMatchingPredicate(ctx context.Context, id RestorableDatabaseAccountId, predicate RestorableGremlinDatabaseGetResultOperationPredicate) (result RestorableGremlinDatabasesListCompleteResult, err error) {
	items := make([]RestorableGremlinDatabaseGetResult, 0)

	resp, err := c.RestorableGremlinDatabasesList(ctx, id)
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

	result = RestorableGremlinDatabasesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
