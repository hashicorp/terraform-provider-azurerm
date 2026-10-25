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

type LoadBalancerProbesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Probe
}

type LoadBalancerProbesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Probe
}

type LoadBalancerProbesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LoadBalancerProbesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LoadBalancerProbesList ...
func (c LoadBalancersClient) LoadBalancerProbesList(ctx context.Context, id ProviderLoadBalancerId) (result LoadBalancerProbesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &LoadBalancerProbesListCustomPager{},
		Path:       fmt.Sprintf("%s/probes", id.ID()),
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
		Values *[]Probe `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LoadBalancerProbesListComplete retrieves all the results into a single object
func (c LoadBalancersClient) LoadBalancerProbesListComplete(ctx context.Context, id ProviderLoadBalancerId) (LoadBalancerProbesListCompleteResult, error) {
	return c.LoadBalancerProbesListCompleteMatchingPredicate(ctx, id, ProbeOperationPredicate{})
}

// LoadBalancerProbesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LoadBalancersClient) LoadBalancerProbesListCompleteMatchingPredicate(ctx context.Context, id ProviderLoadBalancerId, predicate ProbeOperationPredicate) (result LoadBalancerProbesListCompleteResult, err error) {
	items := make([]Probe, 0)

	resp, err := c.LoadBalancerProbesList(ctx, id)
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

	result = LoadBalancerProbesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
