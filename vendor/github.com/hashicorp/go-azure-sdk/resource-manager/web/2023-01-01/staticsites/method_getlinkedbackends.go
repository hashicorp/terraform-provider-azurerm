package staticsites

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetLinkedBackendsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StaticSiteLinkedBackendARMResource
}

type GetLinkedBackendsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StaticSiteLinkedBackendARMResource
}

type GetLinkedBackendsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetLinkedBackendsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetLinkedBackends ...
func (c StaticSitesClient) GetLinkedBackends(ctx context.Context, id StaticSiteId) (result GetLinkedBackendsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetLinkedBackendsCustomPager{},
		Path:       fmt.Sprintf("%s/linkedBackends", id.ID()),
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
		Values *[]StaticSiteLinkedBackendARMResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetLinkedBackendsComplete retrieves all the results into a single object
func (c StaticSitesClient) GetLinkedBackendsComplete(ctx context.Context, id StaticSiteId) (GetLinkedBackendsCompleteResult, error) {
	return c.GetLinkedBackendsCompleteMatchingPredicate(ctx, id, StaticSiteLinkedBackendARMResourceOperationPredicate{})
}

// GetLinkedBackendsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) GetLinkedBackendsCompleteMatchingPredicate(ctx context.Context, id StaticSiteId, predicate StaticSiteLinkedBackendARMResourceOperationPredicate) (result GetLinkedBackendsCompleteResult, err error) {
	items := make([]StaticSiteLinkedBackendARMResource, 0)

	resp, err := c.GetLinkedBackends(ctx, id)
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

	result = GetLinkedBackendsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
