package batchaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchAccountListOutboundNetworkDependenciesEndpointsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]OutboundEnvironmentEndpoint
}

type BatchAccountListOutboundNetworkDependenciesEndpointsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []OutboundEnvironmentEndpoint
}

type BatchAccountListOutboundNetworkDependenciesEndpointsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *BatchAccountListOutboundNetworkDependenciesEndpointsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// BatchAccountListOutboundNetworkDependenciesEndpoints ...
func (c BatchAccountsClient) BatchAccountListOutboundNetworkDependenciesEndpoints(ctx context.Context, id BatchAccountId) (result BatchAccountListOutboundNetworkDependenciesEndpointsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &BatchAccountListOutboundNetworkDependenciesEndpointsCustomPager{},
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

// BatchAccountListOutboundNetworkDependenciesEndpointsComplete retrieves all the results into a single object
func (c BatchAccountsClient) BatchAccountListOutboundNetworkDependenciesEndpointsComplete(ctx context.Context, id BatchAccountId) (BatchAccountListOutboundNetworkDependenciesEndpointsCompleteResult, error) {
	return c.BatchAccountListOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate(ctx, id, OutboundEnvironmentEndpointOperationPredicate{})
}

// BatchAccountListOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c BatchAccountsClient) BatchAccountListOutboundNetworkDependenciesEndpointsCompleteMatchingPredicate(ctx context.Context, id BatchAccountId, predicate OutboundEnvironmentEndpointOperationPredicate) (result BatchAccountListOutboundNetworkDependenciesEndpointsCompleteResult, err error) {
	items := make([]OutboundEnvironmentEndpoint, 0)

	resp, err := c.BatchAccountListOutboundNetworkDependenciesEndpoints(ctx, id)
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

	result = BatchAccountListOutboundNetworkDependenciesEndpointsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
