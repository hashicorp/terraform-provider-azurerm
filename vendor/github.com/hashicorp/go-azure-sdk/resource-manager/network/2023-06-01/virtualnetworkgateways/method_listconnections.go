package virtualnetworkgateways

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListConnectionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualNetworkGatewayConnectionListEntity
}

type ListConnectionsCompleteResult struct {
	Items []VirtualNetworkGatewayConnectionListEntity
}

// ListConnections ...
func (c VirtualNetworkGatewaysClient) ListConnections(ctx context.Context, id VirtualNetworkGatewayId) (result ListConnectionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/connections", id.ID()),
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
		Values *[]VirtualNetworkGatewayConnectionListEntity `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListConnectionsComplete retrieves all the results into a single object
func (c VirtualNetworkGatewaysClient) ListConnectionsComplete(ctx context.Context, id VirtualNetworkGatewayId) (ListConnectionsCompleteResult, error) {
	return c.ListConnectionsCompleteMatchingPredicate(ctx, id, VirtualNetworkGatewayConnectionListEntityOperationPredicate{})
}

// ListConnectionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualNetworkGatewaysClient) ListConnectionsCompleteMatchingPredicate(ctx context.Context, id VirtualNetworkGatewayId, predicate VirtualNetworkGatewayConnectionListEntityOperationPredicate) (result ListConnectionsCompleteResult, err error) {
	items := make([]VirtualNetworkGatewayConnectionListEntity, 0)

	resp, err := c.ListConnections(ctx, id)
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

	result = ListConnectionsCompleteResult{
		Items: items,
	}
	return
}
