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

type MongoDBResourcesListMongoUserDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MongoUserDefinitionGetResults
}

type MongoDBResourcesListMongoUserDefinitionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []MongoUserDefinitionGetResults
}

type MongoDBResourcesListMongoUserDefinitionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MongoDBResourcesListMongoUserDefinitionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MongoDBResourcesListMongoUserDefinitions ...
func (c OpenapisClient) MongoDBResourcesListMongoUserDefinitions(ctx context.Context, id DatabaseAccountId) (result MongoDBResourcesListMongoUserDefinitionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &MongoDBResourcesListMongoUserDefinitionsCustomPager{},
		Path:       fmt.Sprintf("%s/mongodbUserDefinitions", id.ID()),
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
		Values *[]MongoUserDefinitionGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MongoDBResourcesListMongoUserDefinitionsComplete retrieves all the results into a single object
func (c OpenapisClient) MongoDBResourcesListMongoUserDefinitionsComplete(ctx context.Context, id DatabaseAccountId) (MongoDBResourcesListMongoUserDefinitionsCompleteResult, error) {
	return c.MongoDBResourcesListMongoUserDefinitionsCompleteMatchingPredicate(ctx, id, MongoUserDefinitionGetResultsOperationPredicate{})
}

// MongoDBResourcesListMongoUserDefinitionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) MongoDBResourcesListMongoUserDefinitionsCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate MongoUserDefinitionGetResultsOperationPredicate) (result MongoDBResourcesListMongoUserDefinitionsCompleteResult, err error) {
	items := make([]MongoUserDefinitionGetResults, 0)

	resp, err := c.MongoDBResourcesListMongoUserDefinitions(ctx, id)
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

	result = MongoDBResourcesListMongoUserDefinitionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
