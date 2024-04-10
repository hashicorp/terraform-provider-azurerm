package resourceproviders

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetUsagesInLocationlistOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CsmUsageQuota
}

type GetUsagesInLocationlistCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CsmUsageQuota
}

// GetUsagesInLocationlist ...
func (c ResourceProvidersClient) GetUsagesInLocationlist(ctx context.Context, id ProviderLocationId) (result GetUsagesInLocationlistOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
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
		Values *[]CsmUsageQuota `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetUsagesInLocationlistComplete retrieves all the results into a single object
func (c ResourceProvidersClient) GetUsagesInLocationlistComplete(ctx context.Context, id ProviderLocationId) (GetUsagesInLocationlistCompleteResult, error) {
	return c.GetUsagesInLocationlistCompleteMatchingPredicate(ctx, id, CsmUsageQuotaOperationPredicate{})
}

// GetUsagesInLocationlistCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceProvidersClient) GetUsagesInLocationlistCompleteMatchingPredicate(ctx context.Context, id ProviderLocationId, predicate CsmUsageQuotaOperationPredicate) (result GetUsagesInLocationlistCompleteResult, err error) {
	items := make([]CsmUsageQuota, 0)

	resp, err := c.GetUsagesInLocationlist(ctx, id)
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

	result = GetUsagesInLocationlistCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
