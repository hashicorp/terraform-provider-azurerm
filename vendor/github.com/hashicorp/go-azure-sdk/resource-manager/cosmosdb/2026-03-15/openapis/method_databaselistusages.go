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

type DatabaseListUsagesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Usage
}

type DatabaseListUsagesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Usage
}

type DatabaseListUsagesOperationOptions struct {
	Filter *string
}

func DefaultDatabaseListUsagesOperationOptions() DatabaseListUsagesOperationOptions {
	return DatabaseListUsagesOperationOptions{}
}

func (o DatabaseListUsagesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DatabaseListUsagesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o DatabaseListUsagesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type DatabaseListUsagesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DatabaseListUsagesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DatabaseListUsages ...
func (c OpenapisClient) DatabaseListUsages(ctx context.Context, id DatabaseId, options DatabaseListUsagesOperationOptions) (result DatabaseListUsagesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &DatabaseListUsagesCustomPager{},
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

// DatabaseListUsagesComplete retrieves all the results into a single object
func (c OpenapisClient) DatabaseListUsagesComplete(ctx context.Context, id DatabaseId, options DatabaseListUsagesOperationOptions) (DatabaseListUsagesCompleteResult, error) {
	return c.DatabaseListUsagesCompleteMatchingPredicate(ctx, id, options, UsageOperationPredicate{})
}

// DatabaseListUsagesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) DatabaseListUsagesCompleteMatchingPredicate(ctx context.Context, id DatabaseId, options DatabaseListUsagesOperationOptions, predicate UsageOperationPredicate) (result DatabaseListUsagesCompleteResult, err error) {
	items := make([]Usage, 0)

	resp, err := c.DatabaseListUsages(ctx, id, options)
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

	result = DatabaseListUsagesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
