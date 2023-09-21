package jobs

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

type ExecutionsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]JobExecution

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ExecutionsListOperationResponse, error)
}

type ExecutionsListCompleteResult struct {
	Items []JobExecution
}

func (r ExecutionsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ExecutionsListOperationResponse) LoadMore(ctx context.Context) (resp ExecutionsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ExecutionsListOperationOptions struct {
	Filter *string
}

func DefaultExecutionsListOperationOptions() ExecutionsListOperationOptions {
	return ExecutionsListOperationOptions{}
}

func (o ExecutionsListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ExecutionsListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// ExecutionsList ...
func (c JobsClient) ExecutionsList(ctx context.Context, id JobId, options ExecutionsListOperationOptions) (resp ExecutionsListOperationResponse, err error) {
	req, err := c.preparerForExecutionsList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "jobs.JobsClient", "ExecutionsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "jobs.JobsClient", "ExecutionsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForExecutionsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "jobs.JobsClient", "ExecutionsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForExecutionsList prepares the ExecutionsList request.
func (c JobsClient) preparerForExecutionsList(ctx context.Context, id JobId, options ExecutionsListOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/executions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForExecutionsListWithNextLink prepares the ExecutionsList request with the given nextLink token.
func (c JobsClient) preparerForExecutionsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForExecutionsList handles the response to the ExecutionsList request. The method always
// closes the http.Response Body.
func (c JobsClient) responderForExecutionsList(resp *http.Response) (result ExecutionsListOperationResponse, err error) {
	type page struct {
		Values   []JobExecution `json:"value"`
		NextLink *string        `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ExecutionsListOperationResponse, err error) {
			req, err := c.preparerForExecutionsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "jobs.JobsClient", "ExecutionsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "jobs.JobsClient", "ExecutionsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForExecutionsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "jobs.JobsClient", "ExecutionsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ExecutionsListComplete retrieves all of the results into a single object
func (c JobsClient) ExecutionsListComplete(ctx context.Context, id JobId, options ExecutionsListOperationOptions) (ExecutionsListCompleteResult, error) {
	return c.ExecutionsListCompleteMatchingPredicate(ctx, id, options, JobExecutionOperationPredicate{})
}

// ExecutionsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c JobsClient) ExecutionsListCompleteMatchingPredicate(ctx context.Context, id JobId, options ExecutionsListOperationOptions, predicate JobExecutionOperationPredicate) (resp ExecutionsListCompleteResult, err error) {
	items := make([]JobExecution, 0)

	page, err := c.ExecutionsList(ctx, id, options)
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

	out := ExecutionsListCompleteResult{
		Items: items,
	}
	return out, nil
}
