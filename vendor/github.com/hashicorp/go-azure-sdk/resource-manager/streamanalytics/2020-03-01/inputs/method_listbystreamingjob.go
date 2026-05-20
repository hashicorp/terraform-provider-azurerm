package inputs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByStreamingJobOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Input
}

type ListByStreamingJobCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Input
}

type ListByStreamingJobOperationOptions struct {
	Select *string
}

func DefaultListByStreamingJobOperationOptions() ListByStreamingJobOperationOptions {
	return ListByStreamingJobOperationOptions{}
}

func (o ListByStreamingJobOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByStreamingJobOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByStreamingJobOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Select != nil {
		out.Append("$select", fmt.Sprintf("%v", *o.Select))
	}
	return &out
}

type ListByStreamingJobCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByStreamingJobCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByStreamingJob ...
func (c InputsClient) ListByStreamingJob(ctx context.Context, id StreamingJobId, options ListByStreamingJobOperationOptions) (result ListByStreamingJobOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByStreamingJobCustomPager{},
		Path:          fmt.Sprintf("%s/inputs", id.ID()),
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
		Values *[]Input `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByStreamingJobComplete retrieves all the results into a single object
func (c InputsClient) ListByStreamingJobComplete(ctx context.Context, id StreamingJobId, options ListByStreamingJobOperationOptions) (ListByStreamingJobCompleteResult, error) {
	return c.ListByStreamingJobCompleteMatchingPredicate(ctx, id, options, InputOperationPredicate{})
}

// ListByStreamingJobCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c InputsClient) ListByStreamingJobCompleteMatchingPredicate(ctx context.Context, id StreamingJobId, options ListByStreamingJobOperationOptions, predicate InputOperationPredicate) (result ListByStreamingJobCompleteResult, err error) {
	items := make([]Input, 0)

	resp, err := c.ListByStreamingJob(ctx, id, options)
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

	result = ListByStreamingJobCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
