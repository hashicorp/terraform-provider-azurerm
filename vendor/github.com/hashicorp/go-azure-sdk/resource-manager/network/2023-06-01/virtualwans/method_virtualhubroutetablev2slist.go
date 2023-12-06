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

type VirtualHubRouteTableV2sListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VirtualHubRouteTableV2
}

type VirtualHubRouteTableV2sListCompleteResult struct {
	Items []VirtualHubRouteTableV2
}

// VirtualHubRouteTableV2sList ...
func (c VirtualWANsClient) VirtualHubRouteTableV2sList(ctx context.Context, id VirtualHubId) (result VirtualHubRouteTableV2sListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/routeTables", id.ID()),
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
		Values *[]VirtualHubRouteTableV2 `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VirtualHubRouteTableV2sListComplete retrieves all the results into a single object
func (c VirtualWANsClient) VirtualHubRouteTableV2sListComplete(ctx context.Context, id VirtualHubId) (VirtualHubRouteTableV2sListCompleteResult, error) {
	return c.VirtualHubRouteTableV2sListCompleteMatchingPredicate(ctx, id, VirtualHubRouteTableV2OperationPredicate{})
}

// VirtualHubRouteTableV2sListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VirtualWANsClient) VirtualHubRouteTableV2sListCompleteMatchingPredicate(ctx context.Context, id VirtualHubId, predicate VirtualHubRouteTableV2OperationPredicate) (result VirtualHubRouteTableV2sListCompleteResult, err error) {
	items := make([]VirtualHubRouteTableV2, 0)

	resp, err := c.VirtualHubRouteTableV2sList(ctx, id)
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

	result = VirtualHubRouteTableV2sListCompleteResult{
		Items: items,
	}
	return
}
