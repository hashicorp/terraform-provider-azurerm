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

type VirtualHubsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualHub
}

type VirtualHubsListByResourceGroupCompleteResult struct {
	Items []VirtualHub
}

// VirtualHubsListByResourceGroup ...
func (c VirtualWANsClient) VirtualHubsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result VirtualHubsListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Network/virtualHubs", id.ID()),
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
		Values *[]VirtualHub `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VirtualHubsListByResourceGroupComplete retrieves all the results into a single object
func (c VirtualWANsClient) VirtualHubsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (VirtualHubsListByResourceGroupCompleteResult, error) {
	return c.VirtualHubsListByResourceGroupCompleteMatchingPredicate(ctx, id, VirtualHubOperationPredicate{})
}

// VirtualHubsListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VirtualHubsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate VirtualHubOperationPredicate) (result VirtualHubsListByResourceGroupCompleteResult, err error) {
	items := make([]VirtualHub, 0)

	resp, err := c.VirtualHubsListByResourceGroup(ctx, id)
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

	result = VirtualHubsListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
