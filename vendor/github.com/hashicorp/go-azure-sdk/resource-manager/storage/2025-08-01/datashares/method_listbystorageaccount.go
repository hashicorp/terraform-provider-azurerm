package datashares

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

type ListByStorageAccountOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DataShare
}

type ListByStorageAccountCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DataShare
}

type ListByStorageAccountCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByStorageAccountCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByStorageAccount ...
func (c DataSharesClient) ListByStorageAccount(ctx context.Context, id commonids.StorageAccountId) (result ListByStorageAccountOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByStorageAccountCustomPager{},
		Path:       fmt.Sprintf("%s/dataShares", id.ID()),
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
		Values *[]DataShare `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByStorageAccountComplete retrieves all the results into a single object
func (c DataSharesClient) ListByStorageAccountComplete(ctx context.Context, id commonids.StorageAccountId) (ListByStorageAccountCompleteResult, error) {
	return c.ListByStorageAccountCompleteMatchingPredicate(ctx, id, DataShareOperationPredicate{})
}

// ListByStorageAccountCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DataSharesClient) ListByStorageAccountCompleteMatchingPredicate(ctx context.Context, id commonids.StorageAccountId, predicate DataShareOperationPredicate) (result ListByStorageAccountCompleteResult, err error) {
	items := make([]DataShare, 0)

	resp, err := c.ListByStorageAccount(ctx, id)
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

	result = ListByStorageAccountCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
