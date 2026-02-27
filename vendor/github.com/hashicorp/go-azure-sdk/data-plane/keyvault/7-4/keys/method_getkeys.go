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

type GetKeysOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]KeyItem
}

type GetKeysCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []KeyItem
}

type GetKeysOperationOptions struct {
	Maxresults *int64
}

func DefaultGetKeysOperationOptions() GetKeysOperationOptions {
	return GetKeysOperationOptions{}
}

func (o GetKeysOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetKeysOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetKeysOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type GetKeysCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetKeysCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetKeys ...
func (c KeysClient) GetKeys(ctx context.Context, options GetKeysOperationOptions) (result GetKeysOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetKeysCustomPager{},
		Path:          "/keys",
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
		Values *[]KeyItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetKeysComplete retrieves all the results into a single object
func (c KeysClient) GetKeysComplete(ctx context.Context, options GetKeysOperationOptions) (GetKeysCompleteResult, error) {
	return c.GetKeysCompleteMatchingPredicate(ctx, options, KeyItemOperationPredicate{})
}

// GetKeysCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c KeysClient) GetKeysCompleteMatchingPredicate(ctx context.Context, options GetKeysOperationOptions, predicate KeyItemOperationPredicate) (result GetKeysCompleteResult, err error) {
	items := make([]KeyItem, 0)

	resp, err := c.GetKeys(ctx, options)
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

	result = GetKeysCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
