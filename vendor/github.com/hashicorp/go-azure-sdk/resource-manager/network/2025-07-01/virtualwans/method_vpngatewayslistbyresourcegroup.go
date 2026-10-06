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

type VpnGatewaysListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VpnGateway
}

type VpnGatewaysListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VpnGateway
}

type VpnGatewaysListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *VpnGatewaysListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// VpnGatewaysListByResourceGroup ...
func (c VirtualWANsClient) VpnGatewaysListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result VpnGatewaysListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &VpnGatewaysListByResourceGroupCustomPager{},
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

// VpnGatewaysListByResourceGroupComplete retrieves all the results into a single object
func (c VirtualWANsClient) VpnGatewaysListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (VpnGatewaysListByResourceGroupCompleteResult, error) {
	return c.VpnGatewaysListByResourceGroupCompleteMatchingPredicate(ctx, id, VpnGatewayOperationPredicate{})
}

// VpnGatewaysListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VpnGatewaysListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate VpnGatewayOperationPredicate) (result VpnGatewaysListByResourceGroupCompleteResult, err error) {
	items := make([]VpnGateway, 0)

	resp, err := c.VpnGatewaysListByResourceGroup(ctx, id)
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

	result = VpnGatewaysListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
