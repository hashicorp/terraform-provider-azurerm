package outputs

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

type ListByStreamingJobOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Output

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByStreamingJobOperationResponse, error)
}

type ListByStreamingJobCompleteResult struct {
	Items []Output
}

func (r ListByStreamingJobOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByStreamingJobOperationResponse) LoadMore(ctx context.Context) (resp ListByStreamingJobOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByStreamingJobOperationOptions struct {
	Select *string
}

func DefaultListByStreamingJobOperationOptions() ListByStreamingJobOperationOptions {
	return ListByStreamingJobOperationOptions{}
}

func (o ListByStreamingJobOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListByStreamingJobOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Select != nil {
		out["$select"] = *o.Select
	}

	return out
}

// ListByStreamingJob ...
func (c OutputsClient) ListByStreamingJob(ctx context.Context, id StreamingJobId, options ListByStreamingJobOperationOptions) (resp ListByStreamingJobOperationResponse, err error) {
	req, err := c.preparerForListByStreamingJob(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "outputs.OutputsClient", "ListByStreamingJob", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "outputs.OutputsClient", "ListByStreamingJob", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByStreamingJob(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "outputs.OutputsClient", "ListByStreamingJob", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByStreamingJob prepares the ListByStreamingJob request.
func (c OutputsClient) preparerForListByStreamingJob(ctx context.Context, id StreamingJobId, options ListByStreamingJobOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/outputs", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByStreamingJobWithNextLink prepares the ListByStreamingJob request with the given nextLink token.
func (c OutputsClient) preparerForListByStreamingJobWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByStreamingJob handles the response to the ListByStreamingJob request. The method always
// closes the http.Response Body.
func (c OutputsClient) responderForListByStreamingJob(resp *http.Response) (result ListByStreamingJobOperationResponse, err error) {
	type page struct {
		Values   []Output `json:"value"`
		NextLink *string  `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByStreamingJobOperationResponse, err error) {
			req, err := c.preparerForListByStreamingJobWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "outputs.OutputsClient", "ListByStreamingJob", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "outputs.OutputsClient", "ListByStreamingJob", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByStreamingJob(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "outputs.OutputsClient", "ListByStreamingJob", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByStreamingJobComplete retrieves all of the results into a single object
func (c OutputsClient) ListByStreamingJobComplete(ctx context.Context, id StreamingJobId, options ListByStreamingJobOperationOptions) (ListByStreamingJobCompleteResult, error) {
	return c.ListByStreamingJobCompleteMatchingPredicate(ctx, id, options, OutputOperationPredicate{})
}

// ListByStreamingJobCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c OutputsClient) ListByStreamingJobCompleteMatchingPredicate(ctx context.Context, id StreamingJobId, options ListByStreamingJobOperationOptions, predicate OutputOperationPredicate) (resp ListByStreamingJobCompleteResult, err error) {
	items := make([]Output, 0)

	page, err := c.ListByStreamingJob(ctx, id, options)
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

	out := ListByStreamingJobCompleteResult{
		Items: items,
	}
	return out, nil
}
