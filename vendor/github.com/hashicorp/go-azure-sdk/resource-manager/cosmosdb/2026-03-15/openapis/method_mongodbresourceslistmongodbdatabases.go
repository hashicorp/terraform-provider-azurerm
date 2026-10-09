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

type MongoDBResourcesListMongoDBDatabasesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MongoDBDatabaseGetResults
}

type MongoDBResourcesListMongoDBDatabasesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []MongoDBDatabaseGetResults
}

type MongoDBResourcesListMongoDBDatabasesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MongoDBResourcesListMongoDBDatabasesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MongoDBResourcesListMongoDBDatabases ...
func (c OpenapisClient) MongoDBResourcesListMongoDBDatabases(ctx context.Context, id DatabaseAccountId) (result MongoDBResourcesListMongoDBDatabasesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &MongoDBResourcesListMongoDBDatabasesCustomPager{},
		Path:       fmt.Sprintf("%s/mongodbDatabases", id.ID()),
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
		Values *[]MongoDBDatabaseGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MongoDBResourcesListMongoDBDatabasesComplete retrieves all the results into a single object
func (c OpenapisClient) MongoDBResourcesListMongoDBDatabasesComplete(ctx context.Context, id DatabaseAccountId) (MongoDBResourcesListMongoDBDatabasesCompleteResult, error) {
	return c.MongoDBResourcesListMongoDBDatabasesCompleteMatchingPredicate(ctx, id, MongoDBDatabaseGetResultsOperationPredicate{})
}

// MongoDBResourcesListMongoDBDatabasesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) MongoDBResourcesListMongoDBDatabasesCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate MongoDBDatabaseGetResultsOperationPredicate) (result MongoDBResourcesListMongoDBDatabasesCompleteResult, err error) {
	items := make([]MongoDBDatabaseGetResults, 0)

	resp, err := c.MongoDBResourcesListMongoDBDatabases(ctx, id)
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

	result = MongoDBResourcesListMongoDBDatabasesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
