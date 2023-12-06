package virtualwans

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RouteMapsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RouteMap
}

type RouteMapsListCompleteResult struct {
	Items []RouteMap
}

// RouteMapsList ...
func (c VirtualWANsClient) RouteMapsList(ctx context.Context, id VirtualHubId) (result RouteMapsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/routeMaps", id.ID()),
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
		Values *[]RouteMap `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// RouteMapsListComplete retrieves all the results into a single object
func (c VirtualWANsClient) RouteMapsListComplete(ctx context.Context, id VirtualHubId) (RouteMapsListCompleteResult, error) {
	return c.RouteMapsListCompleteMatchingPredicate(ctx, id, RouteMapOperationPredicate{})
}

// RouteMapsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) RouteMapsListCompleteMatchingPredicate(ctx context.Context, id VirtualHubId, predicate RouteMapOperationPredicate) (result RouteMapsListCompleteResult, err error) {
	items := make([]RouteMap, 0)

	resp, err := c.RouteMapsList(ctx, id)
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

	result = RouteMapsListCompleteResult{
		Items: items,
	}
	return
}
