package appplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayRouteConfigsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GatewayRouteConfigResource
}

type GatewayRouteConfigsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GatewayRouteConfigResource
}

// GatewayRouteConfigsList ...
func (c AppPlatformClient) GatewayRouteConfigsList(ctx context.Context, id GatewayId) (result GatewayRouteConfigsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/routeConfigs", id.ID()),
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
		Values *[]GatewayRouteConfigResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GatewayRouteConfigsListComplete retrieves all the results into a single object
func (c AppPlatformClient) GatewayRouteConfigsListComplete(ctx context.Context, id GatewayId) (GatewayRouteConfigsListCompleteResult, error) {
	return c.GatewayRouteConfigsListCompleteMatchingPredicate(ctx, id, GatewayRouteConfigResourceOperationPredicate{})
}

// GatewayRouteConfigsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) GatewayRouteConfigsListCompleteMatchingPredicate(ctx context.Context, id GatewayId, predicate GatewayRouteConfigResourceOperationPredicate) (result GatewayRouteConfigsListCompleteResult, err error) {
	items := make([]GatewayRouteConfigResource, 0)

	resp, err := c.GatewayRouteConfigsList(ctx, id)
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

	result = GatewayRouteConfigsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
