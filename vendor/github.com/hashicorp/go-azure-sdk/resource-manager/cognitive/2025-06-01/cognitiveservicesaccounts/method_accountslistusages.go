package cognitiveservicesaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountsListUsagesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Usage
}

type AccountsListUsagesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Usage
}

type AccountsListUsagesOperationOptions struct {
	Filter *string
}

func DefaultAccountsListUsagesOperationOptions() AccountsListUsagesOperationOptions {
	return AccountsListUsagesOperationOptions{}
}

func (o AccountsListUsagesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AccountsListUsagesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AccountsListUsagesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type AccountsListUsagesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AccountsListUsagesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AccountsListUsages ...
func (c CognitiveServicesAccountsClient) AccountsListUsages(ctx context.Context, id AccountId, options AccountsListUsagesOperationOptions) (result AccountsListUsagesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &AccountsListUsagesCustomPager{},
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

// AccountsListUsagesComplete retrieves all the results into a single object
func (c CognitiveServicesAccountsClient) AccountsListUsagesComplete(ctx context.Context, id AccountId, options AccountsListUsagesOperationOptions) (AccountsListUsagesCompleteResult, error) {
	return c.AccountsListUsagesCompleteMatchingPredicate(ctx, id, options, UsageOperationPredicate{})
}

// AccountsListUsagesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CognitiveServicesAccountsClient) AccountsListUsagesCompleteMatchingPredicate(ctx context.Context, id AccountId, options AccountsListUsagesOperationOptions, predicate UsageOperationPredicate) (result AccountsListUsagesCompleteResult, err error) {
	items := make([]Usage, 0)

	resp, err := c.AccountsListUsages(ctx, id, options)
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

	result = AccountsListUsagesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
