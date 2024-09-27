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

type ListConfigurationSnapshotInfoSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SiteConfigurationSnapshotInfo
}

type ListConfigurationSnapshotInfoSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SiteConfigurationSnapshotInfo
}

type ListConfigurationSnapshotInfoSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListConfigurationSnapshotInfoSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListConfigurationSnapshotInfoSlot ...
func (c WebAppsClient) ListConfigurationSnapshotInfoSlot(ctx context.Context, id SlotId) (result ListConfigurationSnapshotInfoSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListConfigurationSnapshotInfoSlotCustomPager{},
		Path:       fmt.Sprintf("%s/config/web/snapshots", id.ID()),
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
		Values *[]SiteConfigurationSnapshotInfo `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListConfigurationSnapshotInfoSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListConfigurationSnapshotInfoSlotComplete(ctx context.Context, id SlotId) (ListConfigurationSnapshotInfoSlotCompleteResult, error) {
	return c.ListConfigurationSnapshotInfoSlotCompleteMatchingPredicate(ctx, id, SiteConfigurationSnapshotInfoOperationPredicate{})
}

// ListConfigurationSnapshotInfoSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListConfigurationSnapshotInfoSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate SiteConfigurationSnapshotInfoOperationPredicate) (result ListConfigurationSnapshotInfoSlotCompleteResult, err error) {
	items := make([]SiteConfigurationSnapshotInfo, 0)

	resp, err := c.ListConfigurationSnapshotInfoSlot(ctx, id)
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

	result = ListConfigurationSnapshotInfoSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
