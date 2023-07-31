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

type VpnSitesListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VpnSite
}

type VpnSitesListByResourceGroupCompleteResult struct {
	Items []VpnSite
}

// VpnSitesListByResourceGroup ...
func (c VirtualWANsClient) VpnSitesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result VpnSitesListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Network/vpnSites", id.ID()),
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
		Values *[]VpnSite `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VpnSitesListByResourceGroupComplete retrieves all the results into a single object
func (c VirtualWANsClient) VpnSitesListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (VpnSitesListByResourceGroupCompleteResult, error) {
	return c.VpnSitesListByResourceGroupCompleteMatchingPredicate(ctx, id, VpnSiteOperationPredicate{})
}

// VpnSitesListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VpnSitesListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate VpnSiteOperationPredicate) (result VpnSitesListByResourceGroupCompleteResult, err error) {
	items := make([]VpnSite, 0)

	resp, err := c.VpnSitesListByResourceGroup(ctx, id)
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

	result = VpnSitesListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
