package hybridrunbookworker

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

type ListByHybridRunbookWorkerGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]HybridRunbookWorker

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByHybridRunbookWorkerGroupOperationResponse, error)
}

type ListByHybridRunbookWorkerGroupCompleteResult struct {
	Items []HybridRunbookWorker
}

func (r ListByHybridRunbookWorkerGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByHybridRunbookWorkerGroupOperationResponse) LoadMore(ctx context.Context) (resp ListByHybridRunbookWorkerGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByHybridRunbookWorkerGroupOperationOptions struct {
	Filter *string
}

func DefaultListByHybridRunbookWorkerGroupOperationOptions() ListByHybridRunbookWorkerGroupOperationOptions {
	return ListByHybridRunbookWorkerGroupOperationOptions{}
}

func (o ListByHybridRunbookWorkerGroupOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListByHybridRunbookWorkerGroupOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// ListByHybridRunbookWorkerGroup ...
func (c HybridRunbookWorkerClient) ListByHybridRunbookWorkerGroup(ctx context.Context, id HybridRunbookWorkerGroupId, options ListByHybridRunbookWorkerGroupOperationOptions) (resp ListByHybridRunbookWorkerGroupOperationResponse, err error) {
	req, err := c.preparerForListByHybridRunbookWorkerGroup(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridrunbookworker.HybridRunbookWorkerClient", "ListByHybridRunbookWorkerGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridrunbookworker.HybridRunbookWorkerClient", "ListByHybridRunbookWorkerGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByHybridRunbookWorkerGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hybridrunbookworker.HybridRunbookWorkerClient", "ListByHybridRunbookWorkerGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByHybridRunbookWorkerGroupComplete retrieves all of the results into a single object
func (c HybridRunbookWorkerClient) ListByHybridRunbookWorkerGroupComplete(ctx context.Context, id HybridRunbookWorkerGroupId, options ListByHybridRunbookWorkerGroupOperationOptions) (ListByHybridRunbookWorkerGroupCompleteResult, error) {
	return c.ListByHybridRunbookWorkerGroupCompleteMatchingPredicate(ctx, id, options, HybridRunbookWorkerOperationPredicate{})
}

// ListByHybridRunbookWorkerGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c HybridRunbookWorkerClient) ListByHybridRunbookWorkerGroupCompleteMatchingPredicate(ctx context.Context, id HybridRunbookWorkerGroupId, options ListByHybridRunbookWorkerGroupOperationOptions, predicate HybridRunbookWorkerOperationPredicate) (resp ListByHybridRunbookWorkerGroupCompleteResult, err error) {
	items := make([]HybridRunbookWorker, 0)

	page, err := c.ListByHybridRunbookWorkerGroup(ctx, id, options)
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

	out := ListByHybridRunbookWorkerGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByHybridRunbookWorkerGroup prepares the ListByHybridRunbookWorkerGroup request.
func (c HybridRunbookWorkerClient) preparerForListByHybridRunbookWorkerGroup(ctx context.Context, id HybridRunbookWorkerGroupId, options ListByHybridRunbookWorkerGroupOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/hybridRunbookWorkers", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByHybridRunbookWorkerGroupWithNextLink prepares the ListByHybridRunbookWorkerGroup request with the given nextLink token.
func (c HybridRunbookWorkerClient) preparerForListByHybridRunbookWorkerGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByHybridRunbookWorkerGroup handles the response to the ListByHybridRunbookWorkerGroup request. The method always
// closes the http.Response Body.
func (c HybridRunbookWorkerClient) responderForListByHybridRunbookWorkerGroup(resp *http.Response) (result ListByHybridRunbookWorkerGroupOperationResponse, err error) {
	type page struct {
		Values   []HybridRunbookWorker `json:"value"`
		NextLink *string               `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByHybridRunbookWorkerGroupOperationResponse, err error) {
			req, err := c.preparerForListByHybridRunbookWorkerGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "hybridrunbookworker.HybridRunbookWorkerClient", "ListByHybridRunbookWorkerGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "hybridrunbookworker.HybridRunbookWorkerClient", "ListByHybridRunbookWorkerGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByHybridRunbookWorkerGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "hybridrunbookworker.HybridRunbookWorkerClient", "ListByHybridRunbookWorkerGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
