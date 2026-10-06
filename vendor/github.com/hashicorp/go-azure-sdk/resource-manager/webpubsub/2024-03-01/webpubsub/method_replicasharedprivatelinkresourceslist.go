package webpubsub

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicaSharedPrivateLinkResourcesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SharedPrivateLinkResource
}

type ReplicaSharedPrivateLinkResourcesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SharedPrivateLinkResource
}

type ReplicaSharedPrivateLinkResourcesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ReplicaSharedPrivateLinkResourcesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ReplicaSharedPrivateLinkResourcesList ...
func (c WebPubSubClient) ReplicaSharedPrivateLinkResourcesList(ctx context.Context, id ReplicaId) (result ReplicaSharedPrivateLinkResourcesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ReplicaSharedPrivateLinkResourcesListCustomPager{},
		Path:       fmt.Sprintf("%s/sharedPrivateLinkResources", id.ID()),
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
		Values *[]SharedPrivateLinkResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ReplicaSharedPrivateLinkResourcesListComplete retrieves all the results into a single object
func (c WebPubSubClient) ReplicaSharedPrivateLinkResourcesListComplete(ctx context.Context, id ReplicaId) (ReplicaSharedPrivateLinkResourcesListCompleteResult, error) {
	return c.ReplicaSharedPrivateLinkResourcesListCompleteMatchingPredicate(ctx, id, SharedPrivateLinkResourceOperationPredicate{})
}

// ReplicaSharedPrivateLinkResourcesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebPubSubClient) ReplicaSharedPrivateLinkResourcesListCompleteMatchingPredicate(ctx context.Context, id ReplicaId, predicate SharedPrivateLinkResourceOperationPredicate) (result ReplicaSharedPrivateLinkResourcesListCompleteResult, err error) {
	items := make([]SharedPrivateLinkResource, 0)

	resp, err := c.ReplicaSharedPrivateLinkResourcesList(ctx, id)
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

	result = ReplicaSharedPrivateLinkResourcesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
