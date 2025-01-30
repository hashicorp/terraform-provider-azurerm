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

type VpnConnectionsListByVpnGatewayOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VpnConnection
}

type VpnConnectionsListByVpnGatewayCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VpnConnection
}

type VpnConnectionsListByVpnGatewayCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *VpnConnectionsListByVpnGatewayCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// VpnConnectionsListByVpnGateway ...
func (c VirtualWANsClient) VpnConnectionsListByVpnGateway(ctx context.Context, id VpnGatewayId) (result VpnConnectionsListByVpnGatewayOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &VpnConnectionsListByVpnGatewayCustomPager{},
		Path:       fmt.Sprintf("%s/vpnConnections", id.ID()),
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
		Values *[]VpnConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VpnConnectionsListByVpnGatewayComplete retrieves all the results into a single object
func (c VirtualWANsClient) VpnConnectionsListByVpnGatewayComplete(ctx context.Context, id VpnGatewayId) (VpnConnectionsListByVpnGatewayCompleteResult, error) {
	return c.VpnConnectionsListByVpnGatewayCompleteMatchingPredicate(ctx, id, VpnConnectionOperationPredicate{})
}

// VpnConnectionsListByVpnGatewayCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VpnConnectionsListByVpnGatewayCompleteMatchingPredicate(ctx context.Context, id VpnGatewayId, predicate VpnConnectionOperationPredicate) (result VpnConnectionsListByVpnGatewayCompleteResult, err error) {
	items := make([]VpnConnection, 0)

	resp, err := c.VpnConnectionsListByVpnGateway(ctx, id)
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

	result = VpnConnectionsListByVpnGatewayCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
