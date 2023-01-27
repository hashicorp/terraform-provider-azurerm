package encodings

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

type JobsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Job

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (JobsListOperationResponse, error)
}

type JobsListCompleteResult struct {
	Items []Job
}

func (r JobsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r JobsListOperationResponse) LoadMore(ctx context.Context) (resp JobsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type JobsListOperationOptions struct {
	Filter  *string
	Orderby *string
}

func DefaultJobsListOperationOptions() JobsListOperationOptions {
	return JobsListOperationOptions{}
}

func (o JobsListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o JobsListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Orderby != nil {
		out["$orderby"] = *o.Orderby
	}

	return out
}

// JobsList ...
func (c EncodingsClient) JobsList(ctx context.Context, id TransformId, options JobsListOperationOptions) (resp JobsListOperationResponse, err error) {
	req, err := c.preparerForJobsList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForJobsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForJobsList prepares the JobsList request.
func (c EncodingsClient) preparerForJobsList(ctx context.Context, id TransformId, options JobsListOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/jobs", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForJobsListWithNextLink prepares the JobsList request with the given nextLink token.
func (c EncodingsClient) preparerForJobsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForJobsList handles the response to the JobsList request. The method always
// closes the http.Response Body.
func (c EncodingsClient) responderForJobsList(resp *http.Response) (result JobsListOperationResponse, err error) {
	type page struct {
		Values   []Job   `json:"value"`
		NextLink *string `json:"@odata.nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result JobsListOperationResponse, err error) {
			req, err := c.preparerForJobsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForJobsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "JobsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// JobsListComplete retrieves all of the results into a single object
func (c EncodingsClient) JobsListComplete(ctx context.Context, id TransformId, options JobsListOperationOptions) (JobsListCompleteResult, error) {
	return c.JobsListCompleteMatchingPredicate(ctx, id, options, JobOperationPredicate{})
}

// JobsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EncodingsClient) JobsListCompleteMatchingPredicate(ctx context.Context, id TransformId, options JobsListOperationOptions, predicate JobOperationPredicate) (resp JobsListCompleteResult, err error) {
	items := make([]Job, 0)

	page, err := c.JobsList(ctx, id, options)
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

	out := JobsListCompleteResult{
		Items: items,
	}
	return out, nil
}
