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

type TableResourcesListTableRoleDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TableRoleDefinitionResource
}

type TableResourcesListTableRoleDefinitionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TableRoleDefinitionResource
}

type TableResourcesListTableRoleDefinitionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *TableResourcesListTableRoleDefinitionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// TableResourcesListTableRoleDefinitions ...
func (c OpenapisClient) TableResourcesListTableRoleDefinitions(ctx context.Context, id DatabaseAccountId) (result TableResourcesListTableRoleDefinitionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &TableResourcesListTableRoleDefinitionsCustomPager{},
		Path:       fmt.Sprintf("%s/tableRoleDefinitions", id.ID()),
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
		Values *[]TableRoleDefinitionResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// TableResourcesListTableRoleDefinitionsComplete retrieves all the results into a single object
func (c OpenapisClient) TableResourcesListTableRoleDefinitionsComplete(ctx context.Context, id DatabaseAccountId) (TableResourcesListTableRoleDefinitionsCompleteResult, error) {
	return c.TableResourcesListTableRoleDefinitionsCompleteMatchingPredicate(ctx, id, TableRoleDefinitionResourceOperationPredicate{})
}

// TableResourcesListTableRoleDefinitionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) TableResourcesListTableRoleDefinitionsCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate TableRoleDefinitionResourceOperationPredicate) (result TableResourcesListTableRoleDefinitionsCompleteResult, err error) {
	items := make([]TableRoleDefinitionResource, 0)

	resp, err := c.TableResourcesListTableRoleDefinitions(ctx, id)
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

	result = TableResourcesListTableRoleDefinitionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
