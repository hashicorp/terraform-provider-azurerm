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

type AssetsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Asset
}

type AssetsListCompleteResult struct {
	Items []Asset
}

type AssetsListOperationOptions struct {
	Filter  *string
	Orderby *string
	Top     *int64
}

func DefaultAssetsListOperationOptions() AssetsListOperationOptions {
	return AssetsListOperationOptions{}
}

func (o AssetsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AssetsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o AssetsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// AssetsList ...
func (c AssetsAndAssetFiltersClient) AssetsList(ctx context.Context, id MediaServiceId, options AssetsListOperationOptions) (result AssetsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/assets", id.ID()),
		OptionsObject: options,
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
		Values *[]Asset `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AssetsListComplete retrieves all the results into a single object
func (c AssetsAndAssetFiltersClient) AssetsListComplete(ctx context.Context, id MediaServiceId, options AssetsListOperationOptions) (AssetsListCompleteResult, error) {
	return c.AssetsListCompleteMatchingPredicate(ctx, id, options, AssetOperationPredicate{})
}

// AssetsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AssetsAndAssetFiltersClient) AssetsListCompleteMatchingPredicate(ctx context.Context, id MediaServiceId, options AssetsListOperationOptions, predicate AssetOperationPredicate) (result AssetsListCompleteResult, err error) {
	items := make([]Asset, 0)

	resp, err := c.AssetsList(ctx, id, options)
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

	result = AssetsListCompleteResult{
		Items: items,
	}
	return
}
