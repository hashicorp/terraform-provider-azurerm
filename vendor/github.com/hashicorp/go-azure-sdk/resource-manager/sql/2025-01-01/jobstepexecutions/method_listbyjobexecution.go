package jobstepexecutions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByJobExecutionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]JobExecution
}

type ListByJobExecutionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []JobExecution
}

type ListByJobExecutionOperationOptions struct {
	CreateTimeMax *string
	CreateTimeMin *string
	EndTimeMax    *string
	EndTimeMin    *string
	IsActive      *bool
	Skip          *int64
	Top           *int64
}

func DefaultListByJobExecutionOperationOptions() ListByJobExecutionOperationOptions {
	return ListByJobExecutionOperationOptions{}
}

func (o ListByJobExecutionOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByJobExecutionOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByJobExecutionOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.CreateTimeMax != nil {
		out.Append("createTimeMax", fmt.Sprintf("%v", *o.CreateTimeMax))
	}
	if o.CreateTimeMin != nil {
		out.Append("createTimeMin", fmt.Sprintf("%v", *o.CreateTimeMin))
	}
	if o.EndTimeMax != nil {
		out.Append("endTimeMax", fmt.Sprintf("%v", *o.EndTimeMax))
	}
	if o.EndTimeMin != nil {
		out.Append("endTimeMin", fmt.Sprintf("%v", *o.EndTimeMin))
	}
	if o.IsActive != nil {
		out.Append("isActive", fmt.Sprintf("%v", *o.IsActive))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListByJobExecutionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByJobExecutionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByJobExecution ...
func (c JobStepExecutionsClient) ListByJobExecution(ctx context.Context, id ExecutionId, options ListByJobExecutionOperationOptions) (result ListByJobExecutionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByJobExecutionCustomPager{},
		Path:          fmt.Sprintf("%s/steps", id.ID()),
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

// ListByJobExecutionComplete retrieves all the results into a single object
func (c JobStepExecutionsClient) ListByJobExecutionComplete(ctx context.Context, id ExecutionId, options ListByJobExecutionOperationOptions) (ListByJobExecutionCompleteResult, error) {
	return c.ListByJobExecutionCompleteMatchingPredicate(ctx, id, options, JobExecutionOperationPredicate{})
}

// ListByJobExecutionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c JobStepExecutionsClient) ListByJobExecutionCompleteMatchingPredicate(ctx context.Context, id ExecutionId, options ListByJobExecutionOperationOptions, predicate JobExecutionOperationPredicate) (result ListByJobExecutionCompleteResult, err error) {
	items := make([]JobExecution, 0)

	resp, err := c.ListByJobExecution(ctx, id, options)
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

	result = ListByJobExecutionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
