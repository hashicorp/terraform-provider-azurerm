package batchaccounts

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

type BatchAccountListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BatchAccount
}

type BatchAccountListByResourceGroupCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BatchAccount
}

type BatchAccountListByResourceGroupCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *BatchAccountListByResourceGroupCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// BatchAccountListByResourceGroup ...
func (c BatchAccountsClient) BatchAccountListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (result BatchAccountListByResourceGroupOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &BatchAccountListByResourceGroupCustomPager{},
		Path:       fmt.Sprintf("%s/providers/Microsoft.Batch/batchAccounts", id.ID()),
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
		Values *[]BatchAccount `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// BatchAccountListByResourceGroupComplete retrieves all the results into a single object
func (c BatchAccountsClient) BatchAccountListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (BatchAccountListByResourceGroupCompleteResult, error) {
	return c.BatchAccountListByResourceGroupCompleteMatchingPredicate(ctx, id, BatchAccountOperationPredicate{})
}

// BatchAccountListByResourceGroupCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c BatchAccountsClient) BatchAccountListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate BatchAccountOperationPredicate) (result BatchAccountListByResourceGroupCompleteResult, err error) {
	items := make([]BatchAccount, 0)

	resp, err := c.BatchAccountListByResourceGroup(ctx, id)
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

	result = BatchAccountListByResourceGroupCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
