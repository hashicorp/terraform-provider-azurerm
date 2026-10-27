package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListSiteContainersSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SiteContainer
}

type ListSiteContainersSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SiteContainer
}

type ListSiteContainersSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListSiteContainersSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListSiteContainersSlot ...
func (c WebAppsClient) ListSiteContainersSlot(ctx context.Context, id SlotId) (result ListSiteContainersSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListSiteContainersSlotCustomPager{},
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

// ListSiteContainersSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListSiteContainersSlotComplete(ctx context.Context, id SlotId) (ListSiteContainersSlotCompleteResult, error) {
	return c.ListSiteContainersSlotCompleteMatchingPredicate(ctx, id, SiteContainerOperationPredicate{})
}

// ListSiteContainersSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListSiteContainersSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate SiteContainerOperationPredicate) (result ListSiteContainersSlotCompleteResult, err error) {
	items := make([]SiteContainer, 0)

	resp, err := c.ListSiteContainersSlot(ctx, id)
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

	result = ListSiteContainersSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
