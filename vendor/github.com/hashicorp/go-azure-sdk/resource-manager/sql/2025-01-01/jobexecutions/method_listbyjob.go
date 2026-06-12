package jobexecutions

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
	Model        *[]JobExecution
}

type ListByJobCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []JobExecution
}

type ListByJobOperationOptions struct {
	CreateTimeMax *string
	CreateTimeMin *string
	EndTimeMax    *string
	EndTimeMin    *string
	IsActive      *bool
	Skip          *int64
	Top           *int64
}

func DefaultListByJobOperationOptions() ListByJobOperationOptions {
	return ListByJobOperationOptions{}
}

func (o ListByJobOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByJobOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByJobOperationOptions) ToQuery() *client.QueryParams {
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
func (c JobExecutionsClient) ListByJob(ctx context.Context, id JobId, options ListByJobOperationOptions) (result ListByJobOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByJobCustomPager{},
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

// ListByJobComplete retrieves all the results into a single object
func (c JobExecutionsClient) ListByJobComplete(ctx context.Context, id JobId, options ListByJobOperationOptions) (ListByJobCompleteResult, error) {
	return c.ListByJobCompleteMatchingPredicate(ctx, id, options, JobExecutionOperationPredicate{})
}

// ListByJobCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c JobExecutionsClient) ListByJobCompleteMatchingPredicate(ctx context.Context, id JobId, options ListByJobOperationOptions, predicate JobExecutionOperationPredicate) (result ListByJobCompleteResult, err error) {
	items := make([]JobExecution, 0)

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
