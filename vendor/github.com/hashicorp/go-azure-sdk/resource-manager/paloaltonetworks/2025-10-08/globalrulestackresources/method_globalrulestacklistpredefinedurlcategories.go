package globalrulestackresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GlobalRulestacklistPredefinedURLCategoriesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PredefinedURLCategory
}

type GlobalRulestacklistPredefinedURLCategoriesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PredefinedURLCategory
}

type GlobalRulestacklistPredefinedURLCategoriesOperationOptions struct {
	Skip *string
	Top  *int64
}

func DefaultGlobalRulestacklistPredefinedURLCategoriesOperationOptions() GlobalRulestacklistPredefinedURLCategoriesOperationOptions {
	return GlobalRulestacklistPredefinedURLCategoriesOperationOptions{}
}

func (o GlobalRulestacklistPredefinedURLCategoriesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GlobalRulestacklistPredefinedURLCategoriesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GlobalRulestacklistPredefinedURLCategoriesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type GlobalRulestacklistPredefinedURLCategoriesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GlobalRulestacklistPredefinedURLCategoriesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GlobalRulestacklistPredefinedURLCategories ...
func (c GlobalRulestackResourcesClient) GlobalRulestacklistPredefinedURLCategories(ctx context.Context, id GlobalRulestackId, options GlobalRulestacklistPredefinedURLCategoriesOperationOptions) (result GlobalRulestacklistPredefinedURLCategoriesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &GlobalRulestacklistPredefinedURLCategoriesCustomPager{},
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

// GlobalRulestacklistPredefinedURLCategoriesComplete retrieves all the results into a single object
func (c GlobalRulestackResourcesClient) GlobalRulestacklistPredefinedURLCategoriesComplete(ctx context.Context, id GlobalRulestackId, options GlobalRulestacklistPredefinedURLCategoriesOperationOptions) (GlobalRulestacklistPredefinedURLCategoriesCompleteResult, error) {
	return c.GlobalRulestacklistPredefinedURLCategoriesCompleteMatchingPredicate(ctx, id, options, PredefinedURLCategoryOperationPredicate{})
}

// GlobalRulestacklistPredefinedURLCategoriesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GlobalRulestackResourcesClient) GlobalRulestacklistPredefinedURLCategoriesCompleteMatchingPredicate(ctx context.Context, id GlobalRulestackId, options GlobalRulestacklistPredefinedURLCategoriesOperationOptions, predicate PredefinedURLCategoryOperationPredicate) (result GlobalRulestacklistPredefinedURLCategoriesCompleteResult, err error) {
	items := make([]PredefinedURLCategory, 0)

	resp, err := c.GlobalRulestacklistPredefinedURLCategories(ctx, id, options)
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

	result = GlobalRulestacklistPredefinedURLCategoriesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
