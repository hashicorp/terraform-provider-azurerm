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

type MongoMIResourcesListMongoMIRoleDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]MongoMIRoleDefinitionResource
}

type MongoMIResourcesListMongoMIRoleDefinitionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []MongoMIRoleDefinitionResource
}

type MongoMIResourcesListMongoMIRoleDefinitionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MongoMIResourcesListMongoMIRoleDefinitionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MongoMIResourcesListMongoMIRoleDefinitions ...
func (c OpenapisClient) MongoMIResourcesListMongoMIRoleDefinitions(ctx context.Context, id DatabaseAccountId) (result MongoMIResourcesListMongoMIRoleDefinitionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &MongoMIResourcesListMongoMIRoleDefinitionsCustomPager{},
		Path:       fmt.Sprintf("%s/mongoMIRoleDefinitions", id.ID()),
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
		Values *[]MongoMIRoleDefinitionResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MongoMIResourcesListMongoMIRoleDefinitionsComplete retrieves all the results into a single object
func (c OpenapisClient) MongoMIResourcesListMongoMIRoleDefinitionsComplete(ctx context.Context, id DatabaseAccountId) (MongoMIResourcesListMongoMIRoleDefinitionsCompleteResult, error) {
	return c.MongoMIResourcesListMongoMIRoleDefinitionsCompleteMatchingPredicate(ctx, id, MongoMIRoleDefinitionResourceOperationPredicate{})
}

// MongoMIResourcesListMongoMIRoleDefinitionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) MongoMIResourcesListMongoMIRoleDefinitionsCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate MongoMIRoleDefinitionResourceOperationPredicate) (result MongoMIResourcesListMongoMIRoleDefinitionsCompleteResult, err error) {
	items := make([]MongoMIRoleDefinitionResource, 0)

	resp, err := c.MongoMIResourcesListMongoMIRoleDefinitions(ctx, id)
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

	result = MongoMIResourcesListMongoMIRoleDefinitionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
