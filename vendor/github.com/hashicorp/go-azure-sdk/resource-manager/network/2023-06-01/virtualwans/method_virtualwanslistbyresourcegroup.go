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

type VirtualWansListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualWAN
}

type VirtualWansListByResourceGroupCompleteResult struct {
	Items []VirtualWAN
}

// VirtualWansListByResourceGroup ...
func (c VirtualWANsClient) VirtualWansListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result VirtualWansListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Network/virtualWans", id.ID()),
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
		Values *[]VirtualWAN `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VirtualWansListByResourceGroupComplete retrieves all the results into a single object
func (c VirtualWANsClient) VirtualWansListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (VirtualWansListByResourceGroupCompleteResult, error) {
	return c.VirtualWansListByResourceGroupCompleteMatchingPredicate(ctx, id, VirtualWANOperationPredicate{})
}

// VirtualWansListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VirtualWansListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate VirtualWANOperationPredicate) (result VirtualWansListByResourceGroupCompleteResult, err error) {
	items := make([]VirtualWAN, 0)

	resp, err := c.VirtualWansListByResourceGroup(ctx, id)
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

	result = VirtualWansListByResourceGroupCompleteResult{
		Items: items,
	}
	return
}
