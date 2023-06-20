package serviceresource

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TasksListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ProjectTask

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (TasksListOperationResponse, error)
}

type TasksListCompleteResult struct {
	Items []ProjectTask
}

func (r TasksListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r TasksListOperationResponse) LoadMore(ctx context.Context) (resp TasksListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type TasksListOperationOptions struct {
	TaskType *string
}

func DefaultTasksListOperationOptions() TasksListOperationOptions {
	return TasksListOperationOptions{}
}

func (o TasksListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o TasksListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.TaskType != nil {
		out["taskType"] = *o.TaskType
	}

	return out
}

// TasksList ...
func (c ServiceResourceClient) TasksList(ctx context.Context, id ProjectId, options TasksListOperationOptions) (resp TasksListOperationResponse, err error) {
	req, err := c.preparerForTasksList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "TasksList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "TasksList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForTasksList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "TasksList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForTasksList prepares the TasksList request.
func (c ServiceResourceClient) preparerForTasksList(ctx context.Context, id ProjectId, options TasksListOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/tasks", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForTasksListWithNextLink prepares the TasksList request with the given nextLink token.
func (c ServiceResourceClient) preparerForTasksListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
	uri, err := url.Parse(nextLink)
	if err != nil {
		return nil, fmt.Errorf("parsing nextLink %q: %+v", nextLink, err)
	}
	queryParameters := map[string]interface{}{}
	for k, v := range uri.Query() {
		if len(v) == 0 {
			continue
		}
		val := v[0]
		val = autorest.Encode("query", val)
		queryParameters[k] = val
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTasksList handles the response to the TasksList request. The method always
// closes the http.Response Body.
func (c ServiceResourceClient) responderForTasksList(resp *http.Response) (result TasksListOperationResponse, err error) {
	type page struct {
		Values   []ProjectTask `json:"value"`
		NextLink *string       `json:"nextLink"`
	}
	var respObj page
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&respObj),
		autorest.ByClosing())
	result.HttpResponse = resp
	result.Model = &respObj.Values
	result.nextLink = respObj.NextLink
	if respObj.NextLink != nil {
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result TasksListOperationResponse, err error) {
			req, err := c.preparerForTasksListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "TasksList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "TasksList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForTasksList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "TasksList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// TasksListComplete retrieves all of the results into a single object
func (c ServiceResourceClient) TasksListComplete(ctx context.Context, id ProjectId, options TasksListOperationOptions) (TasksListCompleteResult, error) {
	return c.TasksListCompleteMatchingPredicate(ctx, id, options, ProjectTaskOperationPredicate{})
}

// TasksListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ServiceResourceClient) TasksListCompleteMatchingPredicate(ctx context.Context, id ProjectId, options TasksListOperationOptions, predicate ProjectTaskOperationPredicate) (resp TasksListCompleteResult, err error) {
	items := make([]ProjectTask, 0)

	page, err := c.TasksList(ctx, id, options)
	if err != nil {
		err = fmt.Errorf("loading the initial page: %+v", err)
		return
	}
	if page.Model != nil {
		for _, v := range *page.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	for page.HasMore() {
		page, err = page.LoadMore(ctx)
		if err != nil {
			err = fmt.Errorf("loading the next page: %+v", err)
			return
		}

		if page.Model != nil {
			for _, v := range *page.Model {
				if predicate.Matches(v) {
					items = append(items, v)
				}
			}
		}
	}

	out := TasksListCompleteResult{
		Items: items,
	}
	return out, nil
}
