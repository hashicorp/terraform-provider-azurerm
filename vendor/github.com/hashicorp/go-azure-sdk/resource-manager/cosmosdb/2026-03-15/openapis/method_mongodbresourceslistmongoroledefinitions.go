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

type MongoDBResourcesListMongoRoleDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MongoRoleDefinitionGetResults
}

type MongoDBResourcesListMongoRoleDefinitionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []MongoRoleDefinitionGetResults
}

type MongoDBResourcesListMongoRoleDefinitionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MongoDBResourcesListMongoRoleDefinitionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MongoDBResourcesListMongoRoleDefinitions ...
func (c OpenapisClient) MongoDBResourcesListMongoRoleDefinitions(ctx context.Context, id DatabaseAccountId) (result MongoDBResourcesListMongoRoleDefinitionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &MongoDBResourcesListMongoRoleDefinitionsCustomPager{},
		Path:       fmt.Sprintf("%s/mongodbRoleDefinitions", id.ID()),
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
		Values *[]MongoRoleDefinitionGetResults `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MongoDBResourcesListMongoRoleDefinitionsComplete retrieves all the results into a single object
func (c OpenapisClient) MongoDBResourcesListMongoRoleDefinitionsComplete(ctx context.Context, id DatabaseAccountId) (MongoDBResourcesListMongoRoleDefinitionsCompleteResult, error) {
	return c.MongoDBResourcesListMongoRoleDefinitionsCompleteMatchingPredicate(ctx, id, MongoRoleDefinitionGetResultsOperationPredicate{})
}

// MongoDBResourcesListMongoRoleDefinitionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) MongoDBResourcesListMongoRoleDefinitionsCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate MongoRoleDefinitionGetResultsOperationPredicate) (result MongoDBResourcesListMongoRoleDefinitionsCompleteResult, err error) {
	items := make([]MongoRoleDefinitionGetResults, 0)

	resp, err := c.MongoDBResourcesListMongoRoleDefinitions(ctx, id)
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

	result = MongoDBResourcesListMongoRoleDefinitionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
