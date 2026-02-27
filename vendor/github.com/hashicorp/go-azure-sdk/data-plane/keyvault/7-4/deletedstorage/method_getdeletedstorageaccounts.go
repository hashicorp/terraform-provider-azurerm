package deletedstorage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetDeletedStorageAccountsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeletedStorageAccountItem
}

type GetDeletedStorageAccountsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeletedStorageAccountItem
}

type GetDeletedStorageAccountsOperationOptions struct {
	Maxresults *int64
}

func DefaultGetDeletedStorageAccountsOperationOptions() GetDeletedStorageAccountsOperationOptions {
	return GetDeletedStorageAccountsOperationOptions{}
}

func (o GetDeletedStorageAccountsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetDeletedStorageAccountsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetDeletedStorageAccountsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type GetDeletedStorageAccountsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetDeletedStorageAccountsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetDeletedStorageAccounts ...
func (c DeletedStorageClient) GetDeletedStorageAccounts(ctx context.Context, options GetDeletedStorageAccountsOperationOptions) (result GetDeletedStorageAccountsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetDeletedStorageAccountsCustomPager{},
		Path:          "/deletedstorage",
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
		Values *[]DeletedStorageAccountItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetDeletedStorageAccountsComplete retrieves all the results into a single object
func (c DeletedStorageClient) GetDeletedStorageAccountsComplete(ctx context.Context, options GetDeletedStorageAccountsOperationOptions) (GetDeletedStorageAccountsCompleteResult, error) {
	return c.GetDeletedStorageAccountsCompleteMatchingPredicate(ctx, options, DeletedStorageAccountItemOperationPredicate{})
}

// GetDeletedStorageAccountsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DeletedStorageClient) GetDeletedStorageAccountsCompleteMatchingPredicate(ctx context.Context, options GetDeletedStorageAccountsOperationOptions, predicate DeletedStorageAccountItemOperationPredicate) (result GetDeletedStorageAccountsCompleteResult, err error) {
	items := make([]DeletedStorageAccountItem, 0)

	resp, err := c.GetDeletedStorageAccounts(ctx, options)
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

	result = GetDeletedStorageAccountsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
