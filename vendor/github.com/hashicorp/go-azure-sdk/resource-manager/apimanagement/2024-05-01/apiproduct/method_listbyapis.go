package apiproduct

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByApisOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProductContract
}

type ListByApisCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProductContract
}

type ListByApisOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultListByApisOperationOptions() ListByApisOperationOptions {
	return ListByApisOperationOptions{}
}

func (o ListByApisOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByApisOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByApisOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListByApisCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByApisCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByApis ...
func (c ApiProductClient) ListByApis(ctx context.Context, id ApiId, options ListByApisOperationOptions) (result ListByApisOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByApisCustomPager{},
		Path:          fmt.Sprintf("%s/products", id.ID()),
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
		Values *[]ProductContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByApisComplete retrieves all the results into a single object
func (c ApiProductClient) ListByApisComplete(ctx context.Context, id ApiId, options ListByApisOperationOptions) (ListByApisCompleteResult, error) {
	return c.ListByApisCompleteMatchingPredicate(ctx, id, options, ProductContractOperationPredicate{})
}

// ListByApisCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiProductClient) ListByApisCompleteMatchingPredicate(ctx context.Context, id ApiId, options ListByApisOperationOptions, predicate ProductContractOperationPredicate) (result ListByApisCompleteResult, err error) {
	items := make([]ProductContract, 0)

	resp, err := c.ListByApis(ctx, id, options)
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

	result = ListByApisCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
