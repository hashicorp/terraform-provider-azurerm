package signalr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListReplicaSkusOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Sku
}

type ListReplicaSkusCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Sku
}

type ListReplicaSkusCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListReplicaSkusCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListReplicaSkus ...
func (c SignalRClient) ListReplicaSkus(ctx context.Context, id ReplicaId) (result ListReplicaSkusOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListReplicaSkusCustomPager{},
		Path:       fmt.Sprintf("%s/skus", id.ID()),
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
		Values *[]Sku `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListReplicaSkusComplete retrieves all the results into a single object
func (c SignalRClient) ListReplicaSkusComplete(ctx context.Context, id ReplicaId) (ListReplicaSkusCompleteResult, error) {
	return c.ListReplicaSkusCompleteMatchingPredicate(ctx, id, SkuOperationPredicate{})
}

// ListReplicaSkusCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SignalRClient) ListReplicaSkusCompleteMatchingPredicate(ctx context.Context, id ReplicaId, predicate SkuOperationPredicate) (result ListReplicaSkusCompleteResult, err error) {
	items := make([]Sku, 0)

	resp, err := c.ListReplicaSkus(ctx, id)
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

	result = ListReplicaSkusCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
