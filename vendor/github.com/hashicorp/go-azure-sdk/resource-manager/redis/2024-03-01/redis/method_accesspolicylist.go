package redis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPolicyListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RedisCacheAccessPolicy
}

type AccessPolicyListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RedisCacheAccessPolicy
}

type AccessPolicyListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AccessPolicyListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AccessPolicyList ...
func (c RedisClient) AccessPolicyList(ctx context.Context, id RediId) (result AccessPolicyListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &AccessPolicyListCustomPager{},
		Path:       fmt.Sprintf("%s/accessPolicies", id.ID()),
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
		Values *[]RedisCacheAccessPolicy `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AccessPolicyListComplete retrieves all the results into a single object
func (c RedisClient) AccessPolicyListComplete(ctx context.Context, id RediId) (AccessPolicyListCompleteResult, error) {
	return c.AccessPolicyListCompleteMatchingPredicate(ctx, id, RedisCacheAccessPolicyOperationPredicate{})
}

// AccessPolicyListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RedisClient) AccessPolicyListCompleteMatchingPredicate(ctx context.Context, id RediId, predicate RedisCacheAccessPolicyOperationPredicate) (result AccessPolicyListCompleteResult, err error) {
	items := make([]RedisCacheAccessPolicy, 0)

	resp, err := c.AccessPolicyList(ctx, id)
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

	result = AccessPolicyListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
