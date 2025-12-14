package jobstream

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByJobOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]JobStream
}

type ListByJobCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []JobStream
}

type ListByJobOperationOptions struct {
	ClientRequestId *string
	Filter          *string
}

func DefaultListByJobOperationOptions() ListByJobOperationOptions {
	return ListByJobOperationOptions{}
}

func (o ListByJobOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.ClientRequestId != nil {
		out.Append("clientRequestId", fmt.Sprintf("%v", *o.ClientRequestId))
	}
	return &out
}

func (o ListByJobOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByJobOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type ListByJobCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByJobCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByJob ...
func (c JobStreamClient) ListByJob(ctx context.Context, id JobId, options ListByJobOperationOptions) (result ListByJobOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByJobCustomPager{},
		Path:          fmt.Sprintf("%s/streams", id.ID()),
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

// ListByJobComplete retrieves all the results into a single object
func (c JobStreamClient) ListByJobComplete(ctx context.Context, id JobId, options ListByJobOperationOptions) (ListByJobCompleteResult, error) {
	return c.ListByJobCompleteMatchingPredicate(ctx, id, options, JobStreamOperationPredicate{})
}

// ListByJobCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c JobStreamClient) ListByJobCompleteMatchingPredicate(ctx context.Context, id JobId, options ListByJobOperationOptions, predicate JobStreamOperationPredicate) (result ListByJobCompleteResult, err error) {
	items := make([]JobStream, 0)

	resp, err := c.ListByJob(ctx, id, options)
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

	result = ListByJobCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
