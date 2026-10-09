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

type DatabaseAccountsListUsagesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Usage
}

type DatabaseAccountsListUsagesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Usage
}

type DatabaseAccountsListUsagesOperationOptions struct {
	Filter *string
}

func DefaultDatabaseAccountsListUsagesOperationOptions() DatabaseAccountsListUsagesOperationOptions {
	return DatabaseAccountsListUsagesOperationOptions{}
}

func (o DatabaseAccountsListUsagesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o DatabaseAccountsListUsagesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o DatabaseAccountsListUsagesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type DatabaseAccountsListUsagesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DatabaseAccountsListUsagesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DatabaseAccountsListUsages ...
func (c OpenapisClient) DatabaseAccountsListUsages(ctx context.Context, id DatabaseAccountId, options DatabaseAccountsListUsagesOperationOptions) (result DatabaseAccountsListUsagesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &DatabaseAccountsListUsagesCustomPager{},
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

// DatabaseAccountsListUsagesComplete retrieves all the results into a single object
func (c OpenapisClient) DatabaseAccountsListUsagesComplete(ctx context.Context, id DatabaseAccountId, options DatabaseAccountsListUsagesOperationOptions) (DatabaseAccountsListUsagesCompleteResult, error) {
	return c.DatabaseAccountsListUsagesCompleteMatchingPredicate(ctx, id, options, UsageOperationPredicate{})
}

// DatabaseAccountsListUsagesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) DatabaseAccountsListUsagesCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, options DatabaseAccountsListUsagesOperationOptions, predicate UsageOperationPredicate) (result DatabaseAccountsListUsagesCompleteResult, err error) {
	items := make([]Usage, 0)

	resp, err := c.DatabaseAccountsListUsages(ctx, id, options)
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

	result = DatabaseAccountsListUsagesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
