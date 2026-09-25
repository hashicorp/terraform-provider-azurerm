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

type GremlinResourcesListGremlinRoleDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GremlinRoleDefinitionResource
}

type GremlinResourcesListGremlinRoleDefinitionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GremlinRoleDefinitionResource
}

type GremlinResourcesListGremlinRoleDefinitionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GremlinResourcesListGremlinRoleDefinitionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GremlinResourcesListGremlinRoleDefinitions ...
func (c OpenapisClient) GremlinResourcesListGremlinRoleDefinitions(ctx context.Context, id DatabaseAccountId) (result GremlinResourcesListGremlinRoleDefinitionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GremlinResourcesListGremlinRoleDefinitionsCustomPager{},
		Path:       fmt.Sprintf("%s/gremlinRoleDefinitions", id.ID()),
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
		Values *[]GremlinRoleDefinitionResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GremlinResourcesListGremlinRoleDefinitionsComplete retrieves all the results into a single object
func (c OpenapisClient) GremlinResourcesListGremlinRoleDefinitionsComplete(ctx context.Context, id DatabaseAccountId) (GremlinResourcesListGremlinRoleDefinitionsCompleteResult, error) {
	return c.GremlinResourcesListGremlinRoleDefinitionsCompleteMatchingPredicate(ctx, id, GremlinRoleDefinitionResourceOperationPredicate{})
}

// GremlinResourcesListGremlinRoleDefinitionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) GremlinResourcesListGremlinRoleDefinitionsCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate GremlinRoleDefinitionResourceOperationPredicate) (result GremlinResourcesListGremlinRoleDefinitionsCompleteResult, err error) {
	items := make([]GremlinRoleDefinitionResource, 0)

	resp, err := c.GremlinResourcesListGremlinRoleDefinitions(ctx, id)
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

	result = GremlinResourcesListGremlinRoleDefinitionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
