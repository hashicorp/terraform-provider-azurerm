package serviceresource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TasksListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProjectTask
}

type TasksListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProjectTask
}

type TasksListOperationOptions struct {
	TaskType *string
}

func DefaultTasksListOperationOptions() TasksListOperationOptions {
	return TasksListOperationOptions{}
}

func (o TasksListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o TasksListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o TasksListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.TaskType != nil {
		out.Append("taskType", fmt.Sprintf("%v", *o.TaskType))
	}
	return &out
}

type TasksListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *TasksListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// TasksList ...
func (c ServiceResourceClient) TasksList(ctx context.Context, id ProjectId, options TasksListOperationOptions) (result TasksListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &TasksListCustomPager{},
		Path:          fmt.Sprintf("%s/tasks", id.ID()),
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
		Values *[]ProjectTask `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// TasksListComplete retrieves all the results into a single object
func (c ServiceResourceClient) TasksListComplete(ctx context.Context, id ProjectId, options TasksListOperationOptions) (TasksListCompleteResult, error) {
	return c.TasksListCompleteMatchingPredicate(ctx, id, options, ProjectTaskOperationPredicate{})
}

// TasksListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ServiceResourceClient) TasksListCompleteMatchingPredicate(ctx context.Context, id ProjectId, options TasksListOperationOptions, predicate ProjectTaskOperationPredicate) (result TasksListCompleteResult, err error) {
	items := make([]ProjectTask, 0)

	resp, err := c.TasksList(ctx, id, options)
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

	result = TasksListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
