package encodings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Job
}

type JobsListCompleteResult struct {
	Items []Job
}

type JobsListOperationOptions struct {
	Filter  *string
	Orderby *string
}

func DefaultJobsListOperationOptions() JobsListOperationOptions {
	return JobsListOperationOptions{}
}

func (o JobsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o JobsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o JobsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	return &out
}

// JobsList ...
func (c EncodingsClient) JobsList(ctx context.Context, id TransformId, options JobsListOperationOptions) (result JobsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/jobs", id.ID()),
		OptionsObject: options,
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
		Values *[]Job `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// JobsListComplete retrieves all the results into a single object
func (c EncodingsClient) JobsListComplete(ctx context.Context, id TransformId, options JobsListOperationOptions) (JobsListCompleteResult, error) {
	return c.JobsListCompleteMatchingPredicate(ctx, id, options, JobOperationPredicate{})
}

// JobsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c EncodingsClient) JobsListCompleteMatchingPredicate(ctx context.Context, id TransformId, options JobsListOperationOptions, predicate JobOperationPredicate) (result JobsListCompleteResult, err error) {
	items := make([]Job, 0)

	resp, err := c.JobsList(ctx, id, options)
	if err != nil {
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

	result = JobsListCompleteResult{
		Items: items,
	}
	return
}
