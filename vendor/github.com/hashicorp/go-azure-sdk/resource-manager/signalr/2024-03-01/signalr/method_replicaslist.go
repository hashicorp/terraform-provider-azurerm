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

type ReplicasListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Replica
}

type ReplicasListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Replica
}

type ReplicasListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ReplicasListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ReplicasList ...
func (c SignalRClient) ReplicasList(ctx context.Context, id SignalRId) (result ReplicasListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ReplicasListCustomPager{},
		Path:       fmt.Sprintf("%s/replicas", id.ID()),
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
		Values *[]Replica `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ReplicasListComplete retrieves all the results into a single object
func (c SignalRClient) ReplicasListComplete(ctx context.Context, id SignalRId) (ReplicasListCompleteResult, error) {
	return c.ReplicasListCompleteMatchingPredicate(ctx, id, ReplicaOperationPredicate{})
}

// ReplicasListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SignalRClient) ReplicasListCompleteMatchingPredicate(ctx context.Context, id SignalRId, predicate ReplicaOperationPredicate) (result ReplicasListCompleteResult, err error) {
	items := make([]Replica, 0)

	resp, err := c.ReplicasList(ctx, id)
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

	result = ReplicasListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
