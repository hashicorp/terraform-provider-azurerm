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

type ListReplicationsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Replication
}

type ListReplicationsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Replication
}

type ListReplicationsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListReplicationsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListReplications ...
func (c VolumesClient) ListReplications(ctx context.Context, id VolumeId) (result ListReplicationsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ListReplicationsCustomPager{},
		Path:       fmt.Sprintf("%s/listReplications", id.ID()),
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
		Values *[]Replication `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListReplicationsComplete retrieves all the results into a single object
func (c VolumesClient) ListReplicationsComplete(ctx context.Context, id VolumeId) (ListReplicationsCompleteResult, error) {
	return c.ListReplicationsCompleteMatchingPredicate(ctx, id, ReplicationOperationPredicate{})
}

// ListReplicationsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VolumesClient) ListReplicationsCompleteMatchingPredicate(ctx context.Context, id VolumeId, predicate ReplicationOperationPredicate) (result ListReplicationsCompleteResult, err error) {
	items := make([]Replication, 0)

	resp, err := c.ListReplications(ctx, id)
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

	result = ListReplicationsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
