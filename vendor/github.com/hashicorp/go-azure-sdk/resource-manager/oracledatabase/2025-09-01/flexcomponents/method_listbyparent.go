package flexcomponents

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByParentOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FlexComponent
}

type ListByParentCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FlexComponent
}

type ListByParentOperationOptions struct {
	Shape *SystemShapes
}

func DefaultListByParentOperationOptions() ListByParentOperationOptions {
	return ListByParentOperationOptions{}
}

func (o ListByParentOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByParentOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByParentOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Shape != nil {
		out.Append("shape", fmt.Sprintf("%v", *o.Shape))
	}
	return &out
}

type ListByParentCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByParentCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByParent ...
func (c FlexComponentsClient) ListByParent(ctx context.Context, id LocationId, options ListByParentOperationOptions) (result ListByParentOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByParentCustomPager{},
		Path:          fmt.Sprintf("%s/flexComponents", id.ID()),
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
		Values *[]FlexComponent `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByParentComplete retrieves all the results into a single object
func (c FlexComponentsClient) ListByParentComplete(ctx context.Context, id LocationId, options ListByParentOperationOptions) (ListByParentCompleteResult, error) {
	return c.ListByParentCompleteMatchingPredicate(ctx, id, options, FlexComponentOperationPredicate{})
}

// ListByParentCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c FlexComponentsClient) ListByParentCompleteMatchingPredicate(ctx context.Context, id LocationId, options ListByParentOperationOptions, predicate FlexComponentOperationPredicate) (result ListByParentCompleteResult, err error) {
	items := make([]FlexComponent, 0)

	resp, err := c.ListByParent(ctx, id, options)
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

	result = ListByParentCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
