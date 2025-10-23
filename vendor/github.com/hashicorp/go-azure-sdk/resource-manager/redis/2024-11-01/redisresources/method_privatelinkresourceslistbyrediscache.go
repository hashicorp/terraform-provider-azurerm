package redisresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkResourcesListByRedisCacheOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PrivateLinkResource
}

type PrivateLinkResourcesListByRedisCacheCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PrivateLinkResource
}

type PrivateLinkResourcesListByRedisCacheCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PrivateLinkResourcesListByRedisCacheCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PrivateLinkResourcesListByRedisCache ...
func (c RedisResourcesClient) PrivateLinkResourcesListByRedisCache(ctx context.Context, id RediId) (result PrivateLinkResourcesListByRedisCacheOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &PrivateLinkResourcesListByRedisCacheCustomPager{},
		Path:       fmt.Sprintf("%s/privateLinkResources", id.ID()),
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
		Values *[]PrivateLinkResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// PrivateLinkResourcesListByRedisCacheComplete retrieves all the results into a single object
func (c RedisResourcesClient) PrivateLinkResourcesListByRedisCacheComplete(ctx context.Context, id RediId) (PrivateLinkResourcesListByRedisCacheCompleteResult, error) {
	return c.PrivateLinkResourcesListByRedisCacheCompleteMatchingPredicate(ctx, id, PrivateLinkResourceOperationPredicate{})
}

// PrivateLinkResourcesListByRedisCacheCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RedisResourcesClient) PrivateLinkResourcesListByRedisCacheCompleteMatchingPredicate(ctx context.Context, id RediId, predicate PrivateLinkResourceOperationPredicate) (result PrivateLinkResourcesListByRedisCacheCompleteResult, err error) {
	items := make([]PrivateLinkResource, 0)

	resp, err := c.PrivateLinkResourcesListByRedisCache(ctx, id)
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

	result = PrivateLinkResourcesListByRedisCacheCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
