package customizationtasks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByCatalogOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CustomizationTask
}

type ListByCatalogCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CustomizationTask
}

type ListByCatalogCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByCatalogCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByCatalog ...
func (c CustomizationTasksClient) ListByCatalog(ctx context.Context, id DevCenterCatalogId) (result ListByCatalogOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByCatalogCustomPager{},
		Path:       fmt.Sprintf("%s/tasks", id.ID()),
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
		Values *[]CustomizationTask `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByCatalogComplete retrieves all the results into a single object
func (c CustomizationTasksClient) ListByCatalogComplete(ctx context.Context, id DevCenterCatalogId) (ListByCatalogCompleteResult, error) {
	return c.ListByCatalogCompleteMatchingPredicate(ctx, id, CustomizationTaskOperationPredicate{})
}

// ListByCatalogCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CustomizationTasksClient) ListByCatalogCompleteMatchingPredicate(ctx context.Context, id DevCenterCatalogId, predicate CustomizationTaskOperationPredicate) (result ListByCatalogCompleteResult, err error) {
	items := make([]CustomizationTask, 0)

	resp, err := c.ListByCatalog(ctx, id)
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

	result = ListByCatalogCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
