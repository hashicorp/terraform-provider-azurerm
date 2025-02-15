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

type ListPredefinedURLCategoriesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PredefinedURLCategory
}

type ListPredefinedURLCategoriesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PredefinedURLCategory
}

type ListPredefinedURLCategoriesOperationOptions struct {
	Skip *string
	Top  *int64
}

func DefaultListPredefinedURLCategoriesOperationOptions() ListPredefinedURLCategoriesOperationOptions {
	return ListPredefinedURLCategoriesOperationOptions{}
}

func (o ListPredefinedURLCategoriesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListPredefinedURLCategoriesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListPredefinedURLCategoriesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListPredefinedURLCategoriesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListPredefinedURLCategoriesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListPredefinedURLCategories ...
func (c GlobalRulestackClient) ListPredefinedURLCategories(ctx context.Context, id GlobalRulestackId, options ListPredefinedURLCategoriesOperationOptions) (result ListPredefinedURLCategoriesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &ListPredefinedURLCategoriesCustomPager{},
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
		Values *[]PredefinedURLCategory `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListPredefinedURLCategoriesComplete retrieves all the results into a single object
func (c GlobalRulestackClient) ListPredefinedURLCategoriesComplete(ctx context.Context, id GlobalRulestackId, options ListPredefinedURLCategoriesOperationOptions) (ListPredefinedURLCategoriesCompleteResult, error) {
	return c.ListPredefinedURLCategoriesCompleteMatchingPredicate(ctx, id, options, PredefinedURLCategoryOperationPredicate{})
}

// ListPredefinedURLCategoriesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GlobalRulestackClient) ListPredefinedURLCategoriesCompleteMatchingPredicate(ctx context.Context, id GlobalRulestackId, options ListPredefinedURLCategoriesOperationOptions, predicate PredefinedURLCategoryOperationPredicate) (result ListPredefinedURLCategoriesCompleteResult, err error) {
	items := make([]PredefinedURLCategory, 0)

	resp, err := c.ListPredefinedURLCategories(ctx, id, options)
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

	result = ListPredefinedURLCategoriesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
