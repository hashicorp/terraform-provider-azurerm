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

type ListInstanceIdentifiersSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]WebSiteInstanceStatus
}

type ListInstanceIdentifiersSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []WebSiteInstanceStatus
}

type ListInstanceIdentifiersSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListInstanceIdentifiersSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListInstanceIdentifiersSlot ...
func (c WebAppsClient) ListInstanceIdentifiersSlot(ctx context.Context, id SlotId) (result ListInstanceIdentifiersSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListInstanceIdentifiersSlotCustomPager{},
		Path:       fmt.Sprintf("%s/instances", id.ID()),
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
		Values *[]WebSiteInstanceStatus `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListInstanceIdentifiersSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListInstanceIdentifiersSlotComplete(ctx context.Context, id SlotId) (ListInstanceIdentifiersSlotCompleteResult, error) {
	return c.ListInstanceIdentifiersSlotCompleteMatchingPredicate(ctx, id, WebSiteInstanceStatusOperationPredicate{})
}

// ListInstanceIdentifiersSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListInstanceIdentifiersSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate WebSiteInstanceStatusOperationPredicate) (result ListInstanceIdentifiersSlotCompleteResult, err error) {
	items := make([]WebSiteInstanceStatus, 0)

	resp, err := c.ListInstanceIdentifiersSlot(ctx, id)
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

	result = ListInstanceIdentifiersSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
