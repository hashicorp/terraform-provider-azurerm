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

type GetInboundNetworkDependenciesEndpointsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]InboundEnvironmentEndpoint
}

type GetInboundNetworkDependenciesEndpointsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []InboundEnvironmentEndpoint
}

type GetInboundNetworkDependenciesEndpointsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetInboundNetworkDependenciesEndpointsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetInboundNetworkDependenciesEndpoints ...
func (c AppServiceEnvironmentsClient) GetInboundNetworkDependenciesEndpoints(ctx context.Context, id commonids.AppServiceEnvironmentId) (result GetInboundNetworkDependenciesEndpointsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetInboundNetworkDependenciesEndpointsCustomPager{},
		Path:       fmt.Sprintf("%s/inboundNetworkDependenciesEndpoints", id.ID()),
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
		Values *[]InboundEnvironmentEndpoint `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetInboundNetworkDependenciesEndpointsComplete retrieves all the results into a single object
func (c AppServiceEnvironmentsClient) GetInboundNetworkDependenciesEndpointsComplete(ctx context.Context, id commonids.AppServiceEnvironmentId) (GetInboundNetworkDependenciesEndpointsCompleteResult, error) {
	return c.GetInboundNetworkDependenciesEndpointsCompleteMatchingPredicate(ctx, id, InboundEnvironmentEndpointOperationPredicate{})
}

// GetInboundNetworkDependenciesEndpointsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppServiceEnvironmentsClient) GetInboundNetworkDependenciesEndpointsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceEnvironmentId, predicate InboundEnvironmentEndpointOperationPredicate) (result GetInboundNetworkDependenciesEndpointsCompleteResult, err error) {
	items := make([]InboundEnvironmentEndpoint, 0)

	resp, err := c.GetInboundNetworkDependenciesEndpoints(ctx, id)
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

	result = GetInboundNetworkDependenciesEndpointsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
