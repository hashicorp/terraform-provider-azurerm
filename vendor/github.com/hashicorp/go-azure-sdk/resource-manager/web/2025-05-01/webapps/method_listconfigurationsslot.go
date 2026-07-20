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

type ListConfigurationsSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SiteConfigResource
}

type ListConfigurationsSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SiteConfigResource
}

type ListConfigurationsSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListConfigurationsSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListConfigurationsSlot ...
func (c WebAppsClient) ListConfigurationsSlot(ctx context.Context, id SlotId) (result ListConfigurationsSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListConfigurationsSlotCustomPager{},
		Path:       fmt.Sprintf("%s/config", id.ID()),
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
		Values *[]SiteConfigResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListConfigurationsSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListConfigurationsSlotComplete(ctx context.Context, id SlotId) (ListConfigurationsSlotCompleteResult, error) {
	return c.ListConfigurationsSlotCompleteMatchingPredicate(ctx, id, SiteConfigResourceOperationPredicate{})
}

// ListConfigurationsSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListConfigurationsSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate SiteConfigResourceOperationPredicate) (result ListConfigurationsSlotCompleteResult, err error) {
	items := make([]SiteConfigResource, 0)

	resp, err := c.ListConfigurationsSlot(ctx, id)
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

	result = ListConfigurationsSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
