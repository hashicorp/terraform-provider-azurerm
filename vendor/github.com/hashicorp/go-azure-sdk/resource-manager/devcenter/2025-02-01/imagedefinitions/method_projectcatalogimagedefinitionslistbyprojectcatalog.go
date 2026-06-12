package imagedefinitions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectCatalogImageDefinitionsListByProjectCatalogOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ImageDefinition
}

type ProjectCatalogImageDefinitionsListByProjectCatalogCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ImageDefinition
}

type ProjectCatalogImageDefinitionsListByProjectCatalogCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ProjectCatalogImageDefinitionsListByProjectCatalogCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ProjectCatalogImageDefinitionsListByProjectCatalog ...
func (c ImageDefinitionsClient) ProjectCatalogImageDefinitionsListByProjectCatalog(ctx context.Context, id CatalogId) (result ProjectCatalogImageDefinitionsListByProjectCatalogOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ProjectCatalogImageDefinitionsListByProjectCatalogCustomPager{},
		Path:       fmt.Sprintf("%s/imageDefinitions", id.ID()),
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
		Values *[]ImageDefinition `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ProjectCatalogImageDefinitionsListByProjectCatalogComplete retrieves all the results into a single object
func (c ImageDefinitionsClient) ProjectCatalogImageDefinitionsListByProjectCatalogComplete(ctx context.Context, id CatalogId) (ProjectCatalogImageDefinitionsListByProjectCatalogCompleteResult, error) {
	return c.ProjectCatalogImageDefinitionsListByProjectCatalogCompleteMatchingPredicate(ctx, id, ImageDefinitionOperationPredicate{})
}

// ProjectCatalogImageDefinitionsListByProjectCatalogCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ImageDefinitionsClient) ProjectCatalogImageDefinitionsListByProjectCatalogCompleteMatchingPredicate(ctx context.Context, id CatalogId, predicate ImageDefinitionOperationPredicate) (result ProjectCatalogImageDefinitionsListByProjectCatalogCompleteResult, err error) {
	items := make([]ImageDefinition, 0)

	resp, err := c.ProjectCatalogImageDefinitionsListByProjectCatalog(ctx, id)
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

	result = ProjectCatalogImageDefinitionsListByProjectCatalogCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
