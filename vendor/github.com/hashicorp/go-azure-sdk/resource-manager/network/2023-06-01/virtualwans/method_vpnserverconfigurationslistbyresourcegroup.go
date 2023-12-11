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

type VpnServerConfigurationsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VpnServerConfiguration
}

type VpnServerConfigurationsListByResourceGroupCompleteResult struct {
	Items []VpnServerConfiguration
}

// VpnServerConfigurationsListByResourceGroup ...
func (c VirtualWANsClient) VpnServerConfigurationsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result VpnServerConfigurationsListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Network/vpnServerConfigurations", id.ID()),
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
		Values *[]VpnServerConfiguration `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VpnServerConfigurationsListByResourceGroupComplete retrieves all the results into a single object
func (c VirtualWANsClient) VpnServerConfigurationsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (VpnServerConfigurationsListByResourceGroupCompleteResult, error) {
	return c.VpnServerConfigurationsListByResourceGroupCompleteMatchingPredicate(ctx, id, VpnServerConfigurationOperationPredicate{})
}

// VpnServerConfigurationsListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VpnServerConfigurationsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate VpnServerConfigurationOperationPredicate) (result VpnServerConfigurationsListByResourceGroupCompleteResult, err error) {
	items := make([]VpnServerConfiguration, 0)

	resp, err := c.VpnServerConfigurationsListByResourceGroup(ctx, id)
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

	result = VpnServerConfigurationsListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
