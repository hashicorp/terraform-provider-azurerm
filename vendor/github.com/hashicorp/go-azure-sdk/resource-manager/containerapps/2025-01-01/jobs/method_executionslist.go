package jobs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExecutionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]JobExecution
}

type ExecutionsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []JobExecution
}

type ExecutionsListOperationOptions struct {
	Filter *string
}

func DefaultExecutionsListOperationOptions() ExecutionsListOperationOptions {
	return ExecutionsListOperationOptions{}
}

func (o ExecutionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ExecutionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ExecutionsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type ExecutionsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ExecutionsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ExecutionsList ...
func (c JobsClient) ExecutionsList(ctx context.Context, id JobId, options ExecutionsListOperationOptions) (result ExecutionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ExecutionsListCustomPager{},
		Path:          fmt.Sprintf("%s/executions", id.ID()),
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
		Values *[]JobExecution `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ExecutionsListComplete retrieves all the results into a single object
func (c JobsClient) ExecutionsListComplete(ctx context.Context, id JobId, options ExecutionsListOperationOptions) (ExecutionsListCompleteResult, error) {
	return c.ExecutionsListCompleteMatchingPredicate(ctx, id, options, JobExecutionOperationPredicate{})
}

// ExecutionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c JobsClient) ExecutionsListCompleteMatchingPredicate(ctx context.Context, id JobId, options ExecutionsListOperationOptions, predicate JobExecutionOperationPredicate) (result ExecutionsListCompleteResult, err error) {
	items := make([]JobExecution, 0)

	resp, err := c.ExecutionsList(ctx, id, options)
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

	result = ExecutionsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
