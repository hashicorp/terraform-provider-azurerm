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

type JobListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CloudJob
}

type JobListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CloudJob
}

type JobListOperationOptions struct {
	ClientRequestId       *string
	Expand                *string
	Filter                *string
	Maxresults            *int64
	OcpDate               *string
	ReturnClientRequestId *bool
	Select                *string
	Timeout               *int64
}

func DefaultJobListOperationOptions() JobListOperationOptions {
	return JobListOperationOptions{}
}

func (o JobListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.ClientRequestId != nil {
		out.Append("client-request-id", fmt.Sprintf("%v", *o.ClientRequestId))
	}
	if o.OcpDate != nil {
		out.Append("ocp-date", fmt.Sprintf("%v", *o.OcpDate))
	}
	if o.ReturnClientRequestId != nil {
		out.Append("return-client-request-id", fmt.Sprintf("%v", *o.ReturnClientRequestId))
	}
	return &out
}

func (o JobListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o JobListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	if o.Select != nil {
		out.Append("$select", fmt.Sprintf("%v", *o.Select))
	}
	if o.Timeout != nil {
		out.Append("timeout", fmt.Sprintf("%v", *o.Timeout))
	}
	return &out
}

type JobListCustomPager struct {
	NextLink *odata.Link `json:"odata.nextLink"`
}

func (p *JobListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// JobList ...
func (c JobsClient) JobList(ctx context.Context, options JobListOperationOptions) (result JobListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &JobListCustomPager{},
		Path:          "/jobs",
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
		Values *[]CloudJob `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// JobListComplete retrieves all the results into a single object
func (c JobsClient) JobListComplete(ctx context.Context, options JobListOperationOptions) (JobListCompleteResult, error) {
	return c.JobListCompleteMatchingPredicate(ctx, options, CloudJobOperationPredicate{})
}

// JobListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c JobsClient) JobListCompleteMatchingPredicate(ctx context.Context, options JobListOperationOptions, predicate CloudJobOperationPredicate) (result JobListCompleteResult, err error) {
	items := make([]CloudJob, 0)

	resp, err := c.JobList(ctx, options)
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

	result = JobListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
