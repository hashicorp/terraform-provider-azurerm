package privatelinkscopedresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByPrivateLinkScopeOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ScopedResource
}

type ListByPrivateLinkScopeCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ScopedResource
}

type ListByPrivateLinkScopeCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByPrivateLinkScopeCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByPrivateLinkScope ...
func (c PrivateLinkScopedResourcesClient) ListByPrivateLinkScope(ctx context.Context, id PrivateLinkScopeId) (result ListByPrivateLinkScopeOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByPrivateLinkScopeCustomPager{},
		Path:       fmt.Sprintf("%s/scopedResources", id.ID()),
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
		Values *[]ScopedResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByPrivateLinkScopeComplete retrieves all the results into a single object
func (c PrivateLinkScopedResourcesClient) ListByPrivateLinkScopeComplete(ctx context.Context, id PrivateLinkScopeId) (ListByPrivateLinkScopeCompleteResult, error) {
	return c.ListByPrivateLinkScopeCompleteMatchingPredicate(ctx, id, ScopedResourceOperationPredicate{})
}

// ListByPrivateLinkScopeCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PrivateLinkScopedResourcesClient) ListByPrivateLinkScopeCompleteMatchingPredicate(ctx context.Context, id PrivateLinkScopeId, predicate ScopedResourceOperationPredicate) (result ListByPrivateLinkScopeCompleteResult, err error) {
	items := make([]ScopedResource, 0)

	resp, err := c.ListByPrivateLinkScope(ctx, id)
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

	result = ListByPrivateLinkScopeCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
