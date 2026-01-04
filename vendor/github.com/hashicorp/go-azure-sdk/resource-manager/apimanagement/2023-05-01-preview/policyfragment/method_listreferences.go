package policyfragment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListReferencesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Resource
}

type ListReferencesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Resource
}

type ListReferencesOperationOptions struct {
	Skip *int64
	Top  *int64
}

func DefaultListReferencesOperationOptions() ListReferencesOperationOptions {
	return ListReferencesOperationOptions{}
}

func (o ListReferencesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListReferencesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListReferencesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListReferencesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListReferencesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListReferences ...
func (c PolicyFragmentClient) ListReferences(ctx context.Context, id PolicyFragmentId, options ListReferencesOperationOptions) (result ListReferencesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &ListReferencesCustomPager{},
		Path:          fmt.Sprintf("%s/listReferences", id.ID()),
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
		Values *[]Resource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListReferencesComplete retrieves all the results into a single object
func (c PolicyFragmentClient) ListReferencesComplete(ctx context.Context, id PolicyFragmentId, options ListReferencesOperationOptions) (ListReferencesCompleteResult, error) {
	return c.ListReferencesCompleteMatchingPredicate(ctx, id, options, ResourceOperationPredicate{})
}

// ListReferencesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PolicyFragmentClient) ListReferencesCompleteMatchingPredicate(ctx context.Context, id PolicyFragmentId, options ListReferencesOperationOptions, predicate ResourceOperationPredicate) (result ListReferencesCompleteResult, err error) {
	items := make([]Resource, 0)

	resp, err := c.ListReferences(ctx, id, options)
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

	result = ListReferencesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
