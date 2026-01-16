package pools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolListByBatchAccountOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Pool
}

type PoolListByBatchAccountCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Pool
}

type PoolListByBatchAccountOperationOptions struct {
	Filter     *string
	Maxresults *int64
	Select     *string
}

func DefaultPoolListByBatchAccountOperationOptions() PoolListByBatchAccountOperationOptions {
	return PoolListByBatchAccountOperationOptions{}
}

func (o PoolListByBatchAccountOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o PoolListByBatchAccountOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o PoolListByBatchAccountOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	if o.Select != nil {
		out.Append("$select", fmt.Sprintf("%v", *o.Select))
	}
	return &out
}

type PoolListByBatchAccountCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *PoolListByBatchAccountCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// PoolListByBatchAccount ...
func (c PoolsClient) PoolListByBatchAccount(ctx context.Context, id BatchAccountId, options PoolListByBatchAccountOperationOptions) (result PoolListByBatchAccountOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &PoolListByBatchAccountCustomPager{},
		Path:          fmt.Sprintf("%s/pools", id.ID()),
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
		Values *[]Pool `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// PoolListByBatchAccountComplete retrieves all the results into a single object
func (c PoolsClient) PoolListByBatchAccountComplete(ctx context.Context, id BatchAccountId, options PoolListByBatchAccountOperationOptions) (PoolListByBatchAccountCompleteResult, error) {
	return c.PoolListByBatchAccountCompleteMatchingPredicate(ctx, id, options, PoolOperationPredicate{})
}

// PoolListByBatchAccountCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PoolsClient) PoolListByBatchAccountCompleteMatchingPredicate(ctx context.Context, id BatchAccountId, options PoolListByBatchAccountOperationOptions, predicate PoolOperationPredicate) (result PoolListByBatchAccountCompleteResult, err error) {
	items := make([]Pool, 0)

	resp, err := c.PoolListByBatchAccount(ctx, id, options)
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

	result = PoolListByBatchAccountCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
