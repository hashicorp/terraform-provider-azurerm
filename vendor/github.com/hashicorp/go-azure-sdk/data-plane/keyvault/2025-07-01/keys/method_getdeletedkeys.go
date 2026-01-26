package keys

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetDeletedKeysOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeletedKeyItem
}

type GetDeletedKeysCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeletedKeyItem
}

type GetDeletedKeysOperationOptions struct {
	Maxresults *int64
}

func DefaultGetDeletedKeysOperationOptions() GetDeletedKeysOperationOptions {
	return GetDeletedKeysOperationOptions{}
}

func (o GetDeletedKeysOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetDeletedKeysOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetDeletedKeysOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type GetDeletedKeysCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetDeletedKeysCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetDeletedKeys ...
func (c KeysClient) GetDeletedKeys(ctx context.Context, options GetDeletedKeysOperationOptions) (result GetDeletedKeysOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetDeletedKeysCustomPager{},
		Path:          "/deletedkeys",
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
		Values *[]DeletedKeyItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetDeletedKeysComplete retrieves all the results into a single object
func (c KeysClient) GetDeletedKeysComplete(ctx context.Context, options GetDeletedKeysOperationOptions) (GetDeletedKeysCompleteResult, error) {
	return c.GetDeletedKeysCompleteMatchingPredicate(ctx, options, DeletedKeyItemOperationPredicate{})
}

// GetDeletedKeysCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c KeysClient) GetDeletedKeysCompleteMatchingPredicate(ctx context.Context, options GetDeletedKeysOperationOptions, predicate DeletedKeyItemOperationPredicate) (result GetDeletedKeysCompleteResult, err error) {
	items := make([]DeletedKeyItem, 0)

	resp, err := c.GetDeletedKeys(ctx, options)
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

	result = GetDeletedKeysCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
