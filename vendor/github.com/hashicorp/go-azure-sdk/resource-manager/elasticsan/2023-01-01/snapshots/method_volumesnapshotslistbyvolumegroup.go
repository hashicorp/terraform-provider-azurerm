package snapshots

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VolumeSnapshotsListByVolumeGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Snapshot
}

type VolumeSnapshotsListByVolumeGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Snapshot
}

type VolumeSnapshotsListByVolumeGroupOperationOptions struct {
	Filter *string
}

func DefaultVolumeSnapshotsListByVolumeGroupOperationOptions() VolumeSnapshotsListByVolumeGroupOperationOptions {
	return VolumeSnapshotsListByVolumeGroupOperationOptions{}
}

func (o VolumeSnapshotsListByVolumeGroupOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o VolumeSnapshotsListByVolumeGroupOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o VolumeSnapshotsListByVolumeGroupOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type VolumeSnapshotsListByVolumeGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *VolumeSnapshotsListByVolumeGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// VolumeSnapshotsListByVolumeGroup ...
func (c SnapshotsClient) VolumeSnapshotsListByVolumeGroup(ctx context.Context, id VolumeGroupId, options VolumeSnapshotsListByVolumeGroupOperationOptions) (result VolumeSnapshotsListByVolumeGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &VolumeSnapshotsListByVolumeGroupCustomPager{},
		Path:          fmt.Sprintf("%s/snapshots", id.ID()),
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

// VolumeSnapshotsListByVolumeGroupComplete retrieves all the results into a single object
func (c SnapshotsClient) VolumeSnapshotsListByVolumeGroupComplete(ctx context.Context, id VolumeGroupId, options VolumeSnapshotsListByVolumeGroupOperationOptions) (VolumeSnapshotsListByVolumeGroupCompleteResult, error) {
	return c.VolumeSnapshotsListByVolumeGroupCompleteMatchingPredicate(ctx, id, options, SnapshotOperationPredicate{})
}

// VolumeSnapshotsListByVolumeGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SnapshotsClient) VolumeSnapshotsListByVolumeGroupCompleteMatchingPredicate(ctx context.Context, id VolumeGroupId, options VolumeSnapshotsListByVolumeGroupOperationOptions, predicate SnapshotOperationPredicate) (result VolumeSnapshotsListByVolumeGroupCompleteResult, err error) {
	items := make([]Snapshot, 0)

	resp, err := c.VolumeSnapshotsListByVolumeGroup(ctx, id, options)
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

	result = VolumeSnapshotsListByVolumeGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
