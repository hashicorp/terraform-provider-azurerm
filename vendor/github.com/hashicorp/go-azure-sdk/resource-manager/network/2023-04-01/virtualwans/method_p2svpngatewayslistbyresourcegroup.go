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

type P2sVpnGatewaysListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]P2SVpnGateway
}

type P2sVpnGatewaysListByResourceGroupCompleteResult struct {
	Items []P2SVpnGateway
}

// P2sVpnGatewaysListByResourceGroup ...
func (c VirtualWANsClient) P2sVpnGatewaysListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result P2sVpnGatewaysListByResourceGroupOperationResponse, err error) {
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

// P2sVpnGatewaysListByResourceGroupComplete retrieves all the results into a single object
func (c VirtualWANsClient) P2sVpnGatewaysListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (P2sVpnGatewaysListByResourceGroupCompleteResult, error) {
	return c.P2sVpnGatewaysListByResourceGroupCompleteMatchingPredicate(ctx, id, P2SVpnGatewayOperationPredicate{})
}

// P2sVpnGatewaysListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) P2sVpnGatewaysListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate P2SVpnGatewayOperationPredicate) (result P2sVpnGatewaysListByResourceGroupCompleteResult, err error) {
	items := make([]P2SVpnGateway, 0)

	resp, err := c.P2sVpnGatewaysListByResourceGroup(ctx, id)
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

	result = P2sVpnGatewaysListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
