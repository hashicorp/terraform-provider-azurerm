package vaults

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UsagesListByVaultsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VaultUsage
}

type UsagesListByVaultsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VaultUsage
}

type UsagesListByVaultsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *UsagesListByVaultsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// UsagesListByVaults ...
func (c VaultsClient) UsagesListByVaults(ctx context.Context, id VaultId) (result UsagesListByVaultsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &UsagesListByVaultsCustomPager{},
		Path:       fmt.Sprintf("%s/usages", id.ID()),
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
		Values *[]VaultUsage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// UsagesListByVaultsComplete retrieves all the results into a single object
func (c VaultsClient) UsagesListByVaultsComplete(ctx context.Context, id VaultId) (UsagesListByVaultsCompleteResult, error) {
	return c.UsagesListByVaultsCompleteMatchingPredicate(ctx, id, VaultUsageOperationPredicate{})
}

// UsagesListByVaultsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c VaultsClient) UsagesListByVaultsCompleteMatchingPredicate(ctx context.Context, id VaultId, predicate VaultUsageOperationPredicate) (result UsagesListByVaultsCompleteResult, err error) {
	items := make([]VaultUsage, 0)

	resp, err := c.UsagesListByVaults(ctx, id)
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

	result = UsagesListByVaultsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
