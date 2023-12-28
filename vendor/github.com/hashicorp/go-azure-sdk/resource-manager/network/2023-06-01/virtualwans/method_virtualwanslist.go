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

type VirtualWansListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualWAN
}

type VirtualWansListCompleteResult struct {
	Items []VirtualWAN
}

// VirtualWansList ...
func (c VirtualWANsClient) VirtualWansList(ctx context.Context, id commonids.SubscriptionId) (result VirtualWansListOperationResponse, err error) {
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

// VirtualWansListComplete retrieves all the results into a single object
func (c VirtualWANsClient) VirtualWansListComplete(ctx context.Context, id commonids.SubscriptionId) (VirtualWansListCompleteResult, error) {
	return c.VirtualWansListCompleteMatchingPredicate(ctx, id, VirtualWANOperationPredicate{})
}

// VirtualWansListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VirtualWansListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate VirtualWANOperationPredicate) (result VirtualWansListCompleteResult, err error) {
	items := make([]VirtualWAN, 0)

	resp, err := c.VirtualWansList(ctx, id)
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

	result = VirtualWansListCompleteResult{
		Items: items,
	}
	return
}
