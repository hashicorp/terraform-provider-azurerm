package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListSiteContainersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SiteContainer
}

type ListSiteContainersCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SiteContainer
}

type ListSiteContainersCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListSiteContainersCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListSiteContainers ...
func (c WebAppsClient) ListSiteContainers(ctx context.Context, id commonids.AppServiceId) (result ListSiteContainersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListSiteContainersCustomPager{},
		Path:       fmt.Sprintf("%s/sitecontainers", id.ID()),
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
		Values *[]SiteContainer `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListSiteContainersComplete retrieves all the results into a single object
func (c WebAppsClient) ListSiteContainersComplete(ctx context.Context, id commonids.AppServiceId) (ListSiteContainersCompleteResult, error) {
	return c.ListSiteContainersCompleteMatchingPredicate(ctx, id, SiteContainerOperationPredicate{})
}

// ListSiteContainersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListSiteContainersCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate SiteContainerOperationPredicate) (result ListSiteContainersCompleteResult, err error) {
	items := make([]SiteContainer, 0)

	resp, err := c.ListSiteContainers(ctx, id)
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

	result = ListSiteContainersCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
