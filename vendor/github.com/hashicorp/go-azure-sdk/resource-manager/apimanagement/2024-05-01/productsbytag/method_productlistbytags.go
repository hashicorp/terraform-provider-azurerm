package productsbytag

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProductListByTagsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TagResourceContract
}

type ProductListByTagsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TagResourceContract
}

type ProductListByTagsOperationOptions struct {
	Filter                   *string
	IncludeNotTaggedProducts *bool
	Skip                     *int64
	Top                      *int64
}

func DefaultProductListByTagsOperationOptions() ProductListByTagsOperationOptions {
	return ProductListByTagsOperationOptions{}
}

func (o ProductListByTagsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ProductListByTagsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ProductListByTagsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.IncludeNotTaggedProducts != nil {
		out.Append("includeNotTaggedProducts", fmt.Sprintf("%v", *o.IncludeNotTaggedProducts))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ProductListByTagsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ProductListByTagsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ProductListByTags ...
func (c ProductsByTagClient) ProductListByTags(ctx context.Context, id ServiceId, options ProductListByTagsOperationOptions) (result ProductListByTagsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ProductListByTagsCustomPager{},
		Path:          fmt.Sprintf("%s/productsByTags", id.ID()),
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
		Values *[]TagResourceContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ProductListByTagsComplete retrieves all the results into a single object
func (c ProductsByTagClient) ProductListByTagsComplete(ctx context.Context, id ServiceId, options ProductListByTagsOperationOptions) (ProductListByTagsCompleteResult, error) {
	return c.ProductListByTagsCompleteMatchingPredicate(ctx, id, options, TagResourceContractOperationPredicate{})
}

// ProductListByTagsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProductsByTagClient) ProductListByTagsCompleteMatchingPredicate(ctx context.Context, id ServiceId, options ProductListByTagsOperationOptions, predicate TagResourceContractOperationPredicate) (result ProductListByTagsCompleteResult, err error) {
	items := make([]TagResourceContract, 0)

	resp, err := c.ProductListByTags(ctx, id, options)
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

	result = ProductListByTagsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
