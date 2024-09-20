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

type ListSnapshotsFromDRSecondarySlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Snapshot
}

type ListSnapshotsFromDRSecondarySlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Snapshot
}

type ListSnapshotsFromDRSecondarySlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListSnapshotsFromDRSecondarySlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListSnapshotsFromDRSecondarySlot ...
func (c WebAppsClient) ListSnapshotsFromDRSecondarySlot(ctx context.Context, id SlotId) (result ListSnapshotsFromDRSecondarySlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListSnapshotsFromDRSecondarySlotCustomPager{},
		Path:       fmt.Sprintf("%s/snapshotsdr", id.ID()),
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
		Values *[]Snapshot `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListSnapshotsFromDRSecondarySlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListSnapshotsFromDRSecondarySlotComplete(ctx context.Context, id SlotId) (ListSnapshotsFromDRSecondarySlotCompleteResult, error) {
	return c.ListSnapshotsFromDRSecondarySlotCompleteMatchingPredicate(ctx, id, SnapshotOperationPredicate{})
}

// ListSnapshotsFromDRSecondarySlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListSnapshotsFromDRSecondarySlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate SnapshotOperationPredicate) (result ListSnapshotsFromDRSecondarySlotCompleteResult, err error) {
	items := make([]Snapshot, 0)

	resp, err := c.ListSnapshotsFromDRSecondarySlot(ctx, id)
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

	result = ListSnapshotsFromDRSecondarySlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
