package volumes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByVolumeGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Volume
}

type ListByVolumeGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Volume
}

type ListByVolumeGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByVolumeGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByVolumeGroup ...
func (c VolumesClient) ListByVolumeGroup(ctx context.Context, id VolumeGroupId) (result ListByVolumeGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByVolumeGroupCustomPager{},
		Path:       fmt.Sprintf("%s/volumes", id.ID()),
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
		Values *[]Volume `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByVolumeGroupComplete retrieves all the results into a single object
func (c VolumesClient) ListByVolumeGroupComplete(ctx context.Context, id VolumeGroupId) (ListByVolumeGroupCompleteResult, error) {
	return c.ListByVolumeGroupCompleteMatchingPredicate(ctx, id, VolumeOperationPredicate{})
}

// ListByVolumeGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VolumesClient) ListByVolumeGroupCompleteMatchingPredicate(ctx context.Context, id VolumeGroupId, predicate VolumeOperationPredicate) (result ListByVolumeGroupCompleteResult, err error) {
	items := make([]Volume, 0)

	resp, err := c.ListByVolumeGroup(ctx, id)
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

	result = ListByVolumeGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
