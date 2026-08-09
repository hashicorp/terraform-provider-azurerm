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

type RestorableMongodbDatabasesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RestorableMongodbDatabaseGetResult
}

type RestorableMongodbDatabasesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RestorableMongodbDatabaseGetResult
}

type RestorableMongodbDatabasesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *RestorableMongodbDatabasesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// RestorableMongodbDatabasesList ...
func (c OpenapisClient) RestorableMongodbDatabasesList(ctx context.Context, id RestorableDatabaseAccountId) (result RestorableMongodbDatabasesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &RestorableMongodbDatabasesListCustomPager{},
		Path:       fmt.Sprintf("%s/restorableMongodbDatabases", id.ID()),
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
		Values *[]RestorableMongodbDatabaseGetResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RestorableMongodbDatabasesListComplete retrieves all the results into a single object
func (c OpenapisClient) RestorableMongodbDatabasesListComplete(ctx context.Context, id RestorableDatabaseAccountId) (RestorableMongodbDatabasesListCompleteResult, error) {
	return c.RestorableMongodbDatabasesListCompleteMatchingPredicate(ctx, id, RestorableMongodbDatabaseGetResultOperationPredicate{})
}

// RestorableMongodbDatabasesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) RestorableMongodbDatabasesListCompleteMatchingPredicate(ctx context.Context, id RestorableDatabaseAccountId, predicate RestorableMongodbDatabaseGetResultOperationPredicate) (result RestorableMongodbDatabasesListCompleteResult, err error) {
	items := make([]RestorableMongodbDatabaseGetResult, 0)

	resp, err := c.RestorableMongodbDatabasesList(ctx, id)
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

	result = RestorableMongodbDatabasesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
