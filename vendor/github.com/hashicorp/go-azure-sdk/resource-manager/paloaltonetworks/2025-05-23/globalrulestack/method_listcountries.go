package globalrulestack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListCountriesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Country
}

type ListCountriesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Country
}

type ListCountriesOperationOptions struct {
	Skip *string
	Top  *int64
}

func DefaultListCountriesOperationOptions() ListCountriesOperationOptions {
	return ListCountriesOperationOptions{}
}

func (o ListCountriesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListCountriesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListCountriesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListCountriesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListCountriesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListCountries ...
func (c GlobalRulestackClient) ListCountries(ctx context.Context, id GlobalRulestackId, options ListCountriesOperationOptions) (result ListCountriesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &ListCountriesCustomPager{},
		Path:          fmt.Sprintf("%s/listCountries", id.ID()),
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
		Values *[]Country `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListCountriesComplete retrieves all the results into a single object
func (c GlobalRulestackClient) ListCountriesComplete(ctx context.Context, id GlobalRulestackId, options ListCountriesOperationOptions) (ListCountriesCompleteResult, error) {
	return c.ListCountriesCompleteMatchingPredicate(ctx, id, options, CountryOperationPredicate{})
}

// ListCountriesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GlobalRulestackClient) ListCountriesCompleteMatchingPredicate(ctx context.Context, id GlobalRulestackId, options ListCountriesOperationOptions, predicate CountryOperationPredicate) (result ListCountriesCompleteResult, err error) {
	items := make([]Country, 0)

	resp, err := c.ListCountries(ctx, id, options)
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

	result = ListCountriesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
