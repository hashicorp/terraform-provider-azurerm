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

type ServiceTasksListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProjectTask
}

type ServiceTasksListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProjectTask
}

type ServiceTasksListOperationOptions struct {
	TaskType *string
}

func DefaultServiceTasksListOperationOptions() ServiceTasksListOperationOptions {
	return ServiceTasksListOperationOptions{}
}

func (o ServiceTasksListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ServiceTasksListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ServiceTasksListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.TaskType != nil {
		out.Append("taskType", fmt.Sprintf("%v", *o.TaskType))
	}
	return &out
}

type ServiceTasksListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ServiceTasksListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ServiceTasksList ...
func (c ServiceResourceClient) ServiceTasksList(ctx context.Context, id ServiceId, options ServiceTasksListOperationOptions) (result ServiceTasksListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ServiceTasksListCustomPager{},
		Path:          fmt.Sprintf("%s/serviceTasks", id.ID()),
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

// ServiceTasksListComplete retrieves all the results into a single object
func (c ServiceResourceClient) ServiceTasksListComplete(ctx context.Context, id ServiceId, options ServiceTasksListOperationOptions) (ServiceTasksListCompleteResult, error) {
	return c.ServiceTasksListCompleteMatchingPredicate(ctx, id, options, ProjectTaskOperationPredicate{})
}

// ServiceTasksListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ServiceResourceClient) ServiceTasksListCompleteMatchingPredicate(ctx context.Context, id ServiceId, options ServiceTasksListOperationOptions, predicate ProjectTaskOperationPredicate) (result ServiceTasksListCompleteResult, err error) {
	items := make([]ProjectTask, 0)

	resp, err := c.ServiceTasksList(ctx, id, options)
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

	result = ServiceTasksListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
