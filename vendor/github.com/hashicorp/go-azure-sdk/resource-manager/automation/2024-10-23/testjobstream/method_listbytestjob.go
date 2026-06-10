package testjobstream

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByTestJobOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]JobStream
}

type ListByTestJobCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []JobStream
}

type ListByTestJobOperationOptions struct {
	Filter *string
}

func DefaultListByTestJobOperationOptions() ListByTestJobOperationOptions {
	return ListByTestJobOperationOptions{}
}

func (o ListByTestJobOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByTestJobOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByTestJobOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type ListByTestJobCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByTestJobCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByTestJob ...
func (c TestJobStreamClient) ListByTestJob(ctx context.Context, id RunbookId, options ListByTestJobOperationOptions) (result ListByTestJobOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByTestJobCustomPager{},
		Path:          fmt.Sprintf("%s/draft/testJob/streams", id.ID()),
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
		Values *[]JobStream `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByTestJobComplete retrieves all the results into a single object
func (c TestJobStreamClient) ListByTestJobComplete(ctx context.Context, id RunbookId, options ListByTestJobOperationOptions) (ListByTestJobCompleteResult, error) {
	return c.ListByTestJobCompleteMatchingPredicate(ctx, id, options, JobStreamOperationPredicate{})
}

// ListByTestJobCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c TestJobStreamClient) ListByTestJobCompleteMatchingPredicate(ctx context.Context, id RunbookId, options ListByTestJobOperationOptions, predicate JobStreamOperationPredicate) (result ListByTestJobCompleteResult, err error) {
	items := make([]JobStream, 0)

	resp, err := c.ListByTestJob(ctx, id, options)
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

	result = ListByTestJobCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
