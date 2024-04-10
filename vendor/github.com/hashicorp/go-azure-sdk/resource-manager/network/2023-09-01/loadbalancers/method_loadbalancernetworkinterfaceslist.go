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

type LoadBalancerNetworkInterfacesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NetworkInterface
}

type LoadBalancerNetworkInterfacesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NetworkInterface
}

// LoadBalancerNetworkInterfacesList ...
func (c LoadBalancersClient) LoadBalancerNetworkInterfacesList(ctx context.Context, id ProviderLoadBalancerId) (result LoadBalancerNetworkInterfacesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/networkInterfaces", id.ID()),
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
		Values *[]NetworkInterface `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LoadBalancerNetworkInterfacesListComplete retrieves all the results into a single object
func (c LoadBalancersClient) LoadBalancerNetworkInterfacesListComplete(ctx context.Context, id ProviderLoadBalancerId) (LoadBalancerNetworkInterfacesListCompleteResult, error) {
	return c.LoadBalancerNetworkInterfacesListCompleteMatchingPredicate(ctx, id, NetworkInterfaceOperationPredicate{})
}

// LoadBalancerNetworkInterfacesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LoadBalancersClient) LoadBalancerNetworkInterfacesListCompleteMatchingPredicate(ctx context.Context, id ProviderLoadBalancerId, predicate NetworkInterfaceOperationPredicate) (result LoadBalancerNetworkInterfacesListCompleteResult, err error) {
	items := make([]NetworkInterface, 0)

	resp, err := c.LoadBalancerNetworkInterfacesList(ctx, id)
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

	result = LoadBalancerNetworkInterfacesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
