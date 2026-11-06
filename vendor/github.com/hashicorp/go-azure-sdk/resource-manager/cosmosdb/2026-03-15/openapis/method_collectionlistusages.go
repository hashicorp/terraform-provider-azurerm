package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CollectionListUsagesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Usage
}

type CollectionListUsagesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Usage
}

type CollectionListUsagesOperationOptions struct {
	Filter *string
}

func DefaultCollectionListUsagesOperationOptions() CollectionListUsagesOperationOptions {
	return CollectionListUsagesOperationOptions{}
}

func (o CollectionListUsagesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CollectionListUsagesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CollectionListUsagesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type CollectionListUsagesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CollectionListUsagesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CollectionListUsages ...
func (c OpenapisClient) CollectionListUsages(ctx context.Context, id CollectionId, options CollectionListUsagesOperationOptions) (result CollectionListUsagesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &CollectionListUsagesCustomPager{},
		Path:          fmt.Sprintf("%s/usages", id.ID()),
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
		Values *[]Usage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CollectionListUsagesComplete retrieves all the results into a single object
func (c OpenapisClient) CollectionListUsagesComplete(ctx context.Context, id CollectionId, options CollectionListUsagesOperationOptions) (CollectionListUsagesCompleteResult, error) {
	return c.CollectionListUsagesCompleteMatchingPredicate(ctx, id, options, UsageOperationPredicate{})
}

// CollectionListUsagesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) CollectionListUsagesCompleteMatchingPredicate(ctx context.Context, id CollectionId, options CollectionListUsagesOperationOptions, predicate UsageOperationPredicate) (result CollectionListUsagesCompleteResult, err error) {
	items := make([]Usage, 0)

	resp, err := c.CollectionListUsages(ctx, id, options)
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

	result = CollectionListUsagesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
