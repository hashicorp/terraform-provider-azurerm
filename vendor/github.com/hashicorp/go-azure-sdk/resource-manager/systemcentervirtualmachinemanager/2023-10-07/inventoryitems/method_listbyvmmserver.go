package inventoryitems

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByVMMServerOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]InventoryItem
}

type ListByVMMServerCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []InventoryItem
}

// ListByVMMServer ...
func (c InventoryItemsClient) ListByVMMServer(ctx context.Context, id VMmServerId) (result ListByVMMServerOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/inventoryItems", id.ID()),
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
		Values *[]InventoryItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByVMMServerComplete retrieves all the results into a single object
func (c InventoryItemsClient) ListByVMMServerComplete(ctx context.Context, id VMmServerId) (ListByVMMServerCompleteResult, error) {
	return c.ListByVMMServerCompleteMatchingPredicate(ctx, id, InventoryItemOperationPredicate{})
}

// ListByVMMServerCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c InventoryItemsClient) ListByVMMServerCompleteMatchingPredicate(ctx context.Context, id VMmServerId, predicate InventoryItemOperationPredicate) (result ListByVMMServerCompleteResult, err error) {
	items := make([]InventoryItem, 0)

	resp, err := c.ListByVMMServer(ctx, id)
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

	result = ListByVMMServerCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
