package loadbalancers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancerBackendAddressPoolsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BackendAddressPool
}

type LoadBalancerBackendAddressPoolsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BackendAddressPool
}

type LoadBalancerBackendAddressPoolsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LoadBalancerBackendAddressPoolsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LoadBalancerBackendAddressPoolsList ...
func (c LoadBalancersClient) LoadBalancerBackendAddressPoolsList(ctx context.Context, id ProviderLoadBalancerId) (result LoadBalancerBackendAddressPoolsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &LoadBalancerBackendAddressPoolsListCustomPager{},
		Path:       fmt.Sprintf("%s/backendAddressPools", id.ID()),
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
		Values *[]BackendAddressPool `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LoadBalancerBackendAddressPoolsListComplete retrieves all the results into a single object
func (c LoadBalancersClient) LoadBalancerBackendAddressPoolsListComplete(ctx context.Context, id ProviderLoadBalancerId) (LoadBalancerBackendAddressPoolsListCompleteResult, error) {
	return c.LoadBalancerBackendAddressPoolsListCompleteMatchingPredicate(ctx, id, BackendAddressPoolOperationPredicate{})
}

// LoadBalancerBackendAddressPoolsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LoadBalancersClient) LoadBalancerBackendAddressPoolsListCompleteMatchingPredicate(ctx context.Context, id ProviderLoadBalancerId, predicate BackendAddressPoolOperationPredicate) (result LoadBalancerBackendAddressPoolsListCompleteResult, err error) {
	items := make([]BackendAddressPool, 0)

	resp, err := c.LoadBalancerBackendAddressPoolsList(ctx, id)
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

	result = LoadBalancerBackendAddressPoolsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
