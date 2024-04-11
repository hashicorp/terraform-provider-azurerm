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

type VpnSiteLinksListByVpnSiteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VpnSiteLink
}

type VpnSiteLinksListByVpnSiteCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VpnSiteLink
}

// VpnSiteLinksListByVpnSite ...
func (c VirtualWANsClient) VpnSiteLinksListByVpnSite(ctx context.Context, id VpnSiteId) (result VpnSiteLinksListByVpnSiteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/vpnSiteLinks", id.ID()),
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
		Values *[]VpnSiteLink `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VpnSiteLinksListByVpnSiteComplete retrieves all the results into a single object
func (c VirtualWANsClient) VpnSiteLinksListByVpnSiteComplete(ctx context.Context, id VpnSiteId) (VpnSiteLinksListByVpnSiteCompleteResult, error) {
	return c.VpnSiteLinksListByVpnSiteCompleteMatchingPredicate(ctx, id, VpnSiteLinkOperationPredicate{})
}

// VpnSiteLinksListByVpnSiteCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VpnSiteLinksListByVpnSiteCompleteMatchingPredicate(ctx context.Context, id VpnSiteId, predicate VpnSiteLinkOperationPredicate) (result VpnSiteLinksListByVpnSiteCompleteResult, err error) {
	items := make([]VpnSiteLink, 0)

	resp, err := c.VpnSiteLinksListByVpnSite(ctx, id)
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

	result = VpnSiteLinksListByVpnSiteCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
