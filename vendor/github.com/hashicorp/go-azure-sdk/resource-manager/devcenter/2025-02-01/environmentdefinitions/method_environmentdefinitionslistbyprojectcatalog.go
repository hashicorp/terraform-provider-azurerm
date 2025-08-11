package environmentdefinitions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentDefinitionsListByProjectCatalogOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]EnvironmentDefinition
}

type EnvironmentDefinitionsListByProjectCatalogCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []EnvironmentDefinition
}

type EnvironmentDefinitionsListByProjectCatalogCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *EnvironmentDefinitionsListByProjectCatalogCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// EnvironmentDefinitionsListByProjectCatalog ...
func (c EnvironmentDefinitionsClient) EnvironmentDefinitionsListByProjectCatalog(ctx context.Context, id CatalogId) (result EnvironmentDefinitionsListByProjectCatalogOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &EnvironmentDefinitionsListByProjectCatalogCustomPager{},
		Path:       fmt.Sprintf("%s/environmentDefinitions", id.ID()),
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
		Values *[]EnvironmentDefinition `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// EnvironmentDefinitionsListByProjectCatalogComplete retrieves all the results into a single object
func (c EnvironmentDefinitionsClient) EnvironmentDefinitionsListByProjectCatalogComplete(ctx context.Context, id CatalogId) (EnvironmentDefinitionsListByProjectCatalogCompleteResult, error) {
	return c.EnvironmentDefinitionsListByProjectCatalogCompleteMatchingPredicate(ctx, id, EnvironmentDefinitionOperationPredicate{})
}

// EnvironmentDefinitionsListByProjectCatalogCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EnvironmentDefinitionsClient) EnvironmentDefinitionsListByProjectCatalogCompleteMatchingPredicate(ctx context.Context, id CatalogId, predicate EnvironmentDefinitionOperationPredicate) (result EnvironmentDefinitionsListByProjectCatalogCompleteResult, err error) {
	items := make([]EnvironmentDefinition, 0)

	resp, err := c.EnvironmentDefinitionsListByProjectCatalog(ctx, id)
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

	result = EnvironmentDefinitionsListByProjectCatalogCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
