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

type ListSnapshotsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Snapshot
}

type ListSnapshotsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Snapshot
}

type ListSnapshotsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListSnapshotsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListSnapshots ...
func (c WebAppsClient) ListSnapshots(ctx context.Context, id commonids.AppServiceId) (result ListSnapshotsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListSnapshotsCustomPager{},
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

// ListSnapshotsComplete retrieves all the results into a single object
func (c WebAppsClient) ListSnapshotsComplete(ctx context.Context, id commonids.AppServiceId) (ListSnapshotsCompleteResult, error) {
	return c.ListSnapshotsCompleteMatchingPredicate(ctx, id, SnapshotOperationPredicate{})
}

// ListSnapshotsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListSnapshotsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate SnapshotOperationPredicate) (result ListSnapshotsCompleteResult, err error) {
	items := make([]Snapshot, 0)

	resp, err := c.ListSnapshots(ctx, id)
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

	result = ListSnapshotsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
