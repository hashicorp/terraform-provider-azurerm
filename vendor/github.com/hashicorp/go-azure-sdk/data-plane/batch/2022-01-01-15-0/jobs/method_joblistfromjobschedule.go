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

type JobListFromJobScheduleOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CloudJob
}

type JobListFromJobScheduleCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CloudJob
}

type JobListFromJobScheduleOperationOptions struct {
	ClientRequestId       *string
	Expand                *string
	Filter                *string
	Maxresults            *int64
	OcpDate               *string
	ReturnClientRequestId *bool
	Select                *string
	Timeout               *int64
}

func DefaultJobListFromJobScheduleOperationOptions() JobListFromJobScheduleOperationOptions {
	return JobListFromJobScheduleOperationOptions{}
}

func (o JobListFromJobScheduleOperationOptions) ToHeaders() *client.Headers {
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

func (o JobListFromJobScheduleOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o JobListFromJobScheduleOperationOptions) ToQuery() *client.QueryParams {
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

type JobListFromJobScheduleCustomPager struct {
	NextLink *odata.Link `json:"odata.nextLink"`
}

func (p *JobListFromJobScheduleCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// JobListFromJobSchedule ...
func (c JobsClient) JobListFromJobSchedule(ctx context.Context, id JobscheduleId, options JobListFromJobScheduleOperationOptions) (result JobListFromJobScheduleOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &JobListFromJobScheduleCustomPager{},
		Path:          fmt.Sprintf("%s/jobs", id.Path()),
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

// JobListFromJobScheduleComplete retrieves all the results into a single object
func (c JobsClient) JobListFromJobScheduleComplete(ctx context.Context, id JobscheduleId, options JobListFromJobScheduleOperationOptions) (JobListFromJobScheduleCompleteResult, error) {
	return c.JobListFromJobScheduleCompleteMatchingPredicate(ctx, id, options, CloudJobOperationPredicate{})
}

// JobListFromJobScheduleCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c JobsClient) JobListFromJobScheduleCompleteMatchingPredicate(ctx context.Context, id JobscheduleId, options JobListFromJobScheduleOperationOptions, predicate CloudJobOperationPredicate) (result JobListFromJobScheduleCompleteResult, err error) {
	items := make([]CloudJob, 0)

	resp, err := c.JobListFromJobSchedule(ctx, id, options)
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

	result = JobListFromJobScheduleCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
