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

type GetStaticSiteBuildsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StaticSiteBuildARMResource
}

type GetStaticSiteBuildsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StaticSiteBuildARMResource
}

type GetStaticSiteBuildsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetStaticSiteBuildsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetStaticSiteBuilds ...
func (c StaticSitesClient) GetStaticSiteBuilds(ctx context.Context, id StaticSiteId) (result GetStaticSiteBuildsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetStaticSiteBuildsCustomPager{},
		Path:       fmt.Sprintf("%s/builds", id.ID()),
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
		Values *[]StaticSiteBuildARMResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetStaticSiteBuildsComplete retrieves all the results into a single object
func (c StaticSitesClient) GetStaticSiteBuildsComplete(ctx context.Context, id StaticSiteId) (GetStaticSiteBuildsCompleteResult, error) {
	return c.GetStaticSiteBuildsCompleteMatchingPredicate(ctx, id, StaticSiteBuildARMResourceOperationPredicate{})
}

// GetStaticSiteBuildsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StaticSitesClient) GetStaticSiteBuildsCompleteMatchingPredicate(ctx context.Context, id StaticSiteId, predicate StaticSiteBuildARMResourceOperationPredicate) (result GetStaticSiteBuildsCompleteResult, err error) {
	items := make([]StaticSiteBuildARMResource, 0)

	resp, err := c.GetStaticSiteBuilds(ctx, id)
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

	result = GetStaticSiteBuildsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
