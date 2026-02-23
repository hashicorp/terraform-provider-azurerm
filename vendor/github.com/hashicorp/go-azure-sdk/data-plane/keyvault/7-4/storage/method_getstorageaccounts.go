package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetStorageAccountsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StorageAccountItem
}

type GetStorageAccountsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StorageAccountItem
}

type GetStorageAccountsOperationOptions struct {
	Maxresults *int64
}

func DefaultGetStorageAccountsOperationOptions() GetStorageAccountsOperationOptions {
	return GetStorageAccountsOperationOptions{}
}

func (o GetStorageAccountsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetStorageAccountsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetStorageAccountsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type GetStorageAccountsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetStorageAccountsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetStorageAccounts ...
func (c StorageClient) GetStorageAccounts(ctx context.Context, options GetStorageAccountsOperationOptions) (result GetStorageAccountsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetStorageAccountsCustomPager{},
		Path:          "/storage",
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
		Values *[]StorageAccountItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetStorageAccountsComplete retrieves all the results into a single object
func (c StorageClient) GetStorageAccountsComplete(ctx context.Context, options GetStorageAccountsOperationOptions) (GetStorageAccountsCompleteResult, error) {
	return c.GetStorageAccountsCompleteMatchingPredicate(ctx, options, StorageAccountItemOperationPredicate{})
}

// GetStorageAccountsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StorageClient) GetStorageAccountsCompleteMatchingPredicate(ctx context.Context, options GetStorageAccountsOperationOptions, predicate StorageAccountItemOperationPredicate) (result GetStorageAccountsCompleteResult, err error) {
	items := make([]StorageAccountItem, 0)

	resp, err := c.GetStorageAccounts(ctx, options)
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

	result = GetStorageAccountsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
