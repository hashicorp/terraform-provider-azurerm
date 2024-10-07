package localrulestacks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListPredefinedUrlCategoriesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PredefinedUrlCategory
}

type ListPredefinedUrlCategoriesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PredefinedUrlCategory
}

type ListPredefinedUrlCategoriesOperationOptions struct {
	Skip *string
	Top  *int64
}

func DefaultListPredefinedUrlCategoriesOperationOptions() ListPredefinedUrlCategoriesOperationOptions {
	return ListPredefinedUrlCategoriesOperationOptions{}
}

func (o ListPredefinedUrlCategoriesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListPredefinedUrlCategoriesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListPredefinedUrlCategoriesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListPredefinedUrlCategoriesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListPredefinedUrlCategoriesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListPredefinedUrlCategories ...
func (c LocalRulestacksClient) ListPredefinedUrlCategories(ctx context.Context, id LocalRulestackId, options ListPredefinedUrlCategoriesOperationOptions) (result ListPredefinedUrlCategoriesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &ListPredefinedUrlCategoriesCustomPager{},
		Path:          fmt.Sprintf("%s/listPredefinedUrlCategories", id.ID()),
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
		Values *[]PredefinedUrlCategory `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListPredefinedUrlCategoriesComplete retrieves all the results into a single object
func (c LocalRulestacksClient) ListPredefinedUrlCategoriesComplete(ctx context.Context, id LocalRulestackId, options ListPredefinedUrlCategoriesOperationOptions) (ListPredefinedUrlCategoriesCompleteResult, error) {
	return c.ListPredefinedUrlCategoriesCompleteMatchingPredicate(ctx, id, options, PredefinedUrlCategoryOperationPredicate{})
}

// ListPredefinedUrlCategoriesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LocalRulestacksClient) ListPredefinedUrlCategoriesCompleteMatchingPredicate(ctx context.Context, id LocalRulestackId, options ListPredefinedUrlCategoriesOperationOptions, predicate PredefinedUrlCategoryOperationPredicate) (result ListPredefinedUrlCategoriesCompleteResult, err error) {
	items := make([]PredefinedUrlCategory, 0)

	resp, err := c.ListPredefinedUrlCategories(ctx, id, options)
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

	result = ListPredefinedUrlCategoriesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
