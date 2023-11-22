package virtualwans

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

type P2sVpnGatewaysListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]P2SVpnGateway
}

type P2sVpnGatewaysListCompleteResult struct {
	Items []P2SVpnGateway
}

// P2sVpnGatewaysList ...
func (c VirtualWANsClient) P2sVpnGatewaysList(ctx context.Context, id commonids.SubscriptionId) (result P2sVpnGatewaysListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Network/p2sVpnGateways", id.ID()),
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
		Values *[]P2SVpnGateway `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// P2sVpnGatewaysListComplete retrieves all the results into a single object
func (c VirtualWANsClient) P2sVpnGatewaysListComplete(ctx context.Context, id commonids.SubscriptionId) (P2sVpnGatewaysListCompleteResult, error) {
	return c.P2sVpnGatewaysListCompleteMatchingPredicate(ctx, id, P2SVpnGatewayOperationPredicate{})
}

// P2sVpnGatewaysListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) P2sVpnGatewaysListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate P2SVpnGatewayOperationPredicate) (result P2sVpnGatewaysListCompleteResult, err error) {
	items := make([]P2SVpnGateway, 0)

	resp, err := c.P2sVpnGatewaysList(ctx, id)
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

	result = P2sVpnGatewaysListCompleteResult{
		Items: items,
	}
	return
}
