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

type GetKeyVersionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]KeyItem
}

type GetKeyVersionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []KeyItem
}

type GetKeyVersionsOperationOptions struct {
	Maxresults *int64
}

func DefaultGetKeyVersionsOperationOptions() GetKeyVersionsOperationOptions {
	return GetKeyVersionsOperationOptions{}
}

func (o GetKeyVersionsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetKeyVersionsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetKeyVersionsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type GetKeyVersionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetKeyVersionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetKeyVersions ...
func (c KeysClient) GetKeyVersions(ctx context.Context, id KeyId, options GetKeyVersionsOperationOptions) (result GetKeyVersionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetKeyVersionsCustomPager{},
		Path:          fmt.Sprintf("%s/versions", id.Path()),
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

// GetKeyVersionsComplete retrieves all the results into a single object
func (c KeysClient) GetKeyVersionsComplete(ctx context.Context, id KeyId, options GetKeyVersionsOperationOptions) (GetKeyVersionsCompleteResult, error) {
	return c.GetKeyVersionsCompleteMatchingPredicate(ctx, id, options, KeyItemOperationPredicate{})
}

// GetKeyVersionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c KeysClient) GetKeyVersionsCompleteMatchingPredicate(ctx context.Context, id KeyId, options GetKeyVersionsOperationOptions, predicate KeyItemOperationPredicate) (result GetKeyVersionsCompleteResult, err error) {
	items := make([]KeyItem, 0)

	resp, err := c.GetKeyVersions(ctx, id, options)
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

	result = GetKeyVersionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
