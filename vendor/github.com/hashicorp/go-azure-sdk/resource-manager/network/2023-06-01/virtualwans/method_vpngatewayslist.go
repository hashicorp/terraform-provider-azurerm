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

type VpnGatewaysListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VpnGateway
}

type VpnGatewaysListCompleteResult struct {
	Items []VpnGateway
}

// VpnGatewaysList ...
func (c VirtualWANsClient) VpnGatewaysList(ctx context.Context, id commonids.SubscriptionId) (result VpnGatewaysListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Network/vpnGateways", id.ID()),
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
		Values *[]VpnGateway `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VpnGatewaysListComplete retrieves all the results into a single object
func (c VirtualWANsClient) VpnGatewaysListComplete(ctx context.Context, id commonids.SubscriptionId) (VpnGatewaysListCompleteResult, error) {
	return c.VpnGatewaysListCompleteMatchingPredicate(ctx, id, VpnGatewayOperationPredicate{})
}

// VpnGatewaysListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VpnGatewaysListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate VpnGatewayOperationPredicate) (result VpnGatewaysListCompleteResult, err error) {
	items := make([]VpnGateway, 0)

	resp, err := c.VpnGatewaysList(ctx, id)
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

	result = VpnGatewaysListCompleteResult{
		Items: items,
	}
	return
}
