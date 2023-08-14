package assetsandassetfilters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetFiltersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AssetFilter
}

type AssetFiltersListCompleteResult struct {
	Items []AssetFilter
}

// AssetFiltersList ...
func (c AssetsAndAssetFiltersClient) AssetFiltersList(ctx context.Context, id AssetId) (result AssetFiltersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/assetFilters", id.ID()),
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
		Values *[]AssetFilter `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AssetFiltersListComplete retrieves all the results into a single object
func (c AssetsAndAssetFiltersClient) AssetFiltersListComplete(ctx context.Context, id AssetId) (AssetFiltersListCompleteResult, error) {
	return c.AssetFiltersListCompleteMatchingPredicate(ctx, id, AssetFilterOperationPredicate{})
}

// AssetFiltersListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AssetsAndAssetFiltersClient) AssetFiltersListCompleteMatchingPredicate(ctx context.Context, id AssetId, predicate AssetFilterOperationPredicate) (result AssetFiltersListCompleteResult, err error) {
	items := make([]AssetFilter, 0)

	resp, err := c.AssetFiltersList(ctx, id)
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

	result = AssetFiltersListCompleteResult{
		Items: items,
	}
	return
}
