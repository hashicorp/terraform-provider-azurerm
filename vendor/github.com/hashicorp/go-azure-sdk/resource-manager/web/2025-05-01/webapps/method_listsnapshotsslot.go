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

type ListSnapshotsSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Snapshot
}

type ListSnapshotsSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Snapshot
}

type ListSnapshotsSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListSnapshotsSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListSnapshotsSlot ...
func (c WebAppsClient) ListSnapshotsSlot(ctx context.Context, id SlotId) (result ListSnapshotsSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListSnapshotsSlotCustomPager{},
		Path:       fmt.Sprintf("%s/snapshots", id.ID()),
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

// ListSnapshotsSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListSnapshotsSlotComplete(ctx context.Context, id SlotId) (ListSnapshotsSlotCompleteResult, error) {
	return c.ListSnapshotsSlotCompleteMatchingPredicate(ctx, id, SnapshotOperationPredicate{})
}

// ListSnapshotsSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListSnapshotsSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate SnapshotOperationPredicate) (result ListSnapshotsSlotCompleteResult, err error) {
	items := make([]Snapshot, 0)

	resp, err := c.ListSnapshotsSlot(ctx, id)
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

	result = ListSnapshotsSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
