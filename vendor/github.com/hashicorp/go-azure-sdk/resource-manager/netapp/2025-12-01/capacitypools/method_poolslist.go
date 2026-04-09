package capacitypools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CapacityPool
}

type PoolsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CapacityPool
}

type PoolsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PoolsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PoolsList ...
func (c CapacityPoolsClient) PoolsList(ctx context.Context, id NetAppAccountId) (result PoolsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &PoolsListCustomPager{},
		Path:       fmt.Sprintf("%s/capacityPools", id.ID()),
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
		Values *[]CapacityPool `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// PoolsListComplete retrieves all the results into a single object
func (c CapacityPoolsClient) PoolsListComplete(ctx context.Context, id NetAppAccountId) (PoolsListCompleteResult, error) {
	return c.PoolsListCompleteMatchingPredicate(ctx, id, CapacityPoolOperationPredicate{})
}

// PoolsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CapacityPoolsClient) PoolsListCompleteMatchingPredicate(ctx context.Context, id NetAppAccountId, predicate CapacityPoolOperationPredicate) (result PoolsListCompleteResult, err error) {
	items := make([]CapacityPool, 0)

	resp, err := c.PoolsList(ctx, id)
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

	result = PoolsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
