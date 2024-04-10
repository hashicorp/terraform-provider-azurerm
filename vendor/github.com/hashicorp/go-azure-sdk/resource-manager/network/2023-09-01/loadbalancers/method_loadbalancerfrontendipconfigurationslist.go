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

type LoadBalancerFrontendIPConfigurationsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FrontendIPConfiguration
}

type LoadBalancerFrontendIPConfigurationsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FrontendIPConfiguration
}

// LoadBalancerFrontendIPConfigurationsList ...
func (c LoadBalancersClient) LoadBalancerFrontendIPConfigurationsList(ctx context.Context, id ProviderLoadBalancerId) (result LoadBalancerFrontendIPConfigurationsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/frontendIPConfigurations", id.ID()),
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
		Values *[]FrontendIPConfiguration `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LoadBalancerFrontendIPConfigurationsListComplete retrieves all the results into a single object
func (c LoadBalancersClient) LoadBalancerFrontendIPConfigurationsListComplete(ctx context.Context, id ProviderLoadBalancerId) (LoadBalancerFrontendIPConfigurationsListCompleteResult, error) {
	return c.LoadBalancerFrontendIPConfigurationsListCompleteMatchingPredicate(ctx, id, FrontendIPConfigurationOperationPredicate{})
}

// LoadBalancerFrontendIPConfigurationsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LoadBalancersClient) LoadBalancerFrontendIPConfigurationsListCompleteMatchingPredicate(ctx context.Context, id ProviderLoadBalancerId, predicate FrontendIPConfigurationOperationPredicate) (result LoadBalancerFrontendIPConfigurationsListCompleteResult, err error) {
	items := make([]FrontendIPConfiguration, 0)

	resp, err := c.LoadBalancerFrontendIPConfigurationsList(ctx, id)
	if err != nil {
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

	result = LoadBalancerFrontendIPConfigurationsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
