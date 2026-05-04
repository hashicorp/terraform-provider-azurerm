package localrulestackresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalRulestackslistPredefinedURLCategoriesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PredefinedURLCategory
}

type LocalRulestackslistPredefinedURLCategoriesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PredefinedURLCategory
}

type LocalRulestackslistPredefinedURLCategoriesOperationOptions struct {
	Skip *string
	Top  *int64
}

func DefaultLocalRulestackslistPredefinedURLCategoriesOperationOptions() LocalRulestackslistPredefinedURLCategoriesOperationOptions {
	return LocalRulestackslistPredefinedURLCategoriesOperationOptions{}
}

func (o LocalRulestackslistPredefinedURLCategoriesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o LocalRulestackslistPredefinedURLCategoriesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o LocalRulestackslistPredefinedURLCategoriesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type LocalRulestackslistPredefinedURLCategoriesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LocalRulestackslistPredefinedURLCategoriesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LocalRulestackslistPredefinedURLCategories ...
func (c LocalRulestackResourcesClient) LocalRulestackslistPredefinedURLCategories(ctx context.Context, id LocalRulestackId, options LocalRulestackslistPredefinedURLCategoriesOperationOptions) (result LocalRulestackslistPredefinedURLCategoriesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &LocalRulestackslistPredefinedURLCategoriesCustomPager{},
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

// LocalRulestackslistPredefinedURLCategoriesComplete retrieves all the results into a single object
func (c LocalRulestackResourcesClient) LocalRulestackslistPredefinedURLCategoriesComplete(ctx context.Context, id LocalRulestackId, options LocalRulestackslistPredefinedURLCategoriesOperationOptions) (LocalRulestackslistPredefinedURLCategoriesCompleteResult, error) {
	return c.LocalRulestackslistPredefinedURLCategoriesCompleteMatchingPredicate(ctx, id, options, PredefinedURLCategoryOperationPredicate{})
}

// LocalRulestackslistPredefinedURLCategoriesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LocalRulestackResourcesClient) LocalRulestackslistPredefinedURLCategoriesCompleteMatchingPredicate(ctx context.Context, id LocalRulestackId, options LocalRulestackslistPredefinedURLCategoriesOperationOptions, predicate PredefinedURLCategoryOperationPredicate) (result LocalRulestackslistPredefinedURLCategoriesCompleteResult, err error) {
	items := make([]PredefinedURLCategory, 0)

	resp, err := c.LocalRulestackslistPredefinedURLCategories(ctx, id, options)
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

	result = LocalRulestackslistPredefinedURLCategoriesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
