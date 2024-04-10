package appserviceenvironments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetOutboundNetworkDependenciesEndpointsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]OutboundEnvironmentEndpoint
}

type GetOutboundNetworkDependenciesEndpointsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []OutboundEnvironmentEndpoint
}

// GetOutboundNetworkDependenciesEndpoints ...
func (c AppServiceEnvironmentsClient) GetOutboundNetworkDependenciesEndpoints(ctx context.Context, id commonids.AppServiceEnvironmentId) (result GetOutboundNetworkDependenciesEndpointsOperationResponse, err error) {
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

// GetOutboundNetworkDependenciesEndpointsComplete retrieves all the results into a single object
func (c AppServiceEnvironmentsClient) GetOutboundNetworkDependenciesEndpointsComplete(ctx context.Context, id commonids.AppServiceEnvironmentId) (GetOutboundNetworkDependenciesEndpointsCompleteResult, error) {
	return c.GetOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate(ctx, id, OutboundEnvironmentEndpointOperationPredicate{})
}

// GetOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppServiceEnvironmentsClient) GetOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceEnvironmentId, predicate OutboundEnvironmentEndpointOperationPredicate) (result GetOutboundNetworkDependenciesEndpointsCompleteResult, err error) {
	items := make([]OutboundEnvironmentEndpoint, 0)

	resp, err := c.GetOutboundNetworkDependenciesEndpoints(ctx, id)
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

	result = GetOutboundNetworkDependenciesEndpointsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
