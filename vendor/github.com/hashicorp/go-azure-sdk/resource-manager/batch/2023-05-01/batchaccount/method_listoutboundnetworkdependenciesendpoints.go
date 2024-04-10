package batchaccount

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListOutboundNetworkDependenciesEndpointsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]OutboundEnvironmentEndpoint
}

type ListOutboundNetworkDependenciesEndpointsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []OutboundEnvironmentEndpoint
}

// ListOutboundNetworkDependenciesEndpoints ...
func (c BatchAccountClient) ListOutboundNetworkDependenciesEndpoints(ctx context.Context, id BatchAccountId) (result ListOutboundNetworkDependenciesEndpointsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/outboundNetworkDependenciesEndpoints", id.ID()),
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
		Values *[]OutboundEnvironmentEndpoint `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListOutboundNetworkDependenciesEndpointsComplete retrieves all the results into a single object
func (c BatchAccountClient) ListOutboundNetworkDependenciesEndpointsComplete(ctx context.Context, id BatchAccountId) (ListOutboundNetworkDependenciesEndpointsCompleteResult, error) {
	return c.ListOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate(ctx, id, OutboundEnvironmentEndpointOperationPredicate{})
}

// ListOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c BatchAccountClient) ListOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate(ctx context.Context, id BatchAccountId, predicate OutboundEnvironmentEndpointOperationPredicate) (result ListOutboundNetworkDependenciesEndpointsCompleteResult, err error) {
	items := make([]OutboundEnvironmentEndpoint, 0)

	resp, err := c.ListOutboundNetworkDependenciesEndpoints(ctx, id)
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

	result = ListOutboundNetworkDependenciesEndpointsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
