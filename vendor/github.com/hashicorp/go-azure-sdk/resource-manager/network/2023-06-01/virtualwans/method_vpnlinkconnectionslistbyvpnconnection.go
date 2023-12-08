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

type VpnLinkConnectionsListByVpnConnectionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VpnSiteLinkConnection
}

type VpnLinkConnectionsListByVpnConnectionCompleteResult struct {
	Items []VpnSiteLinkConnection
}

// VpnLinkConnectionsListByVpnConnection ...
func (c VirtualWANsClient) VpnLinkConnectionsListByVpnConnection(ctx context.Context, id commonids.VPNConnectionId) (result VpnLinkConnectionsListByVpnConnectionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/vpnLinkConnections", id.ID()),
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
		Values *[]VpnSiteLinkConnection `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VpnLinkConnectionsListByVpnConnectionComplete retrieves all the results into a single object
func (c VirtualWANsClient) VpnLinkConnectionsListByVpnConnectionComplete(ctx context.Context, id commonids.VPNConnectionId) (VpnLinkConnectionsListByVpnConnectionCompleteResult, error) {
	return c.VpnLinkConnectionsListByVpnConnectionCompleteMatchingPredicate(ctx, id, VpnSiteLinkConnectionOperationPredicate{})
}

// VpnLinkConnectionsListByVpnConnectionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VpnLinkConnectionsListByVpnConnectionCompleteMatchingPredicate(ctx context.Context, id commonids.VPNConnectionId, predicate VpnSiteLinkConnectionOperationPredicate) (result VpnLinkConnectionsListByVpnConnectionCompleteResult, err error) {
	items := make([]VpnSiteLinkConnection, 0)

	resp, err := c.VpnLinkConnectionsListByVpnConnection(ctx, id)
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

	result = VpnLinkConnectionsListByVpnConnectionCompleteResult{
		Items: items,
	}
	return
}
