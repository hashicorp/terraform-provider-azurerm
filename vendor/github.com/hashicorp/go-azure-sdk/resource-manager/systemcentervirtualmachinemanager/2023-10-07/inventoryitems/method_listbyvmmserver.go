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

type ListByVMmServerOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]InventoryItem
}

type ListByVMmServerCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []InventoryItem
}

type ListByVMmServerCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByVMmServerCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByVMmServer ...
func (c InventoryItemsClient) ListByVMmServer(ctx context.Context, id VMmServerId) (result ListByVMmServerOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByVMmServerCustomPager{},
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

// ListByVMmServerComplete retrieves all the results into a single object
func (c InventoryItemsClient) ListByVMmServerComplete(ctx context.Context, id VMmServerId) (ListByVMmServerCompleteResult, error) {
	return c.ListByVMmServerCompleteMatchingPredicate(ctx, id, InventoryItemOperationPredicate{})
}

// ListByVMmServerCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c InventoryItemsClient) ListByVMmServerCompleteMatchingPredicate(ctx context.Context, id VMmServerId, predicate InventoryItemOperationPredicate) (result ListByVMmServerCompleteResult, err error) {
	items := make([]InventoryItem, 0)

	resp, err := c.ListByVMmServer(ctx, id)
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

	result = ListByVMmServerCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
