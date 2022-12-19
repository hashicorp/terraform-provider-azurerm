package streamingpoliciesandstreaminglocators

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

type StreamingLocatorsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]StreamingLocator

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (StreamingLocatorsListOperationResponse, error)
}

type StreamingLocatorsListCompleteResult struct {
	Items []StreamingLocator
}

func (r StreamingLocatorsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r StreamingLocatorsListOperationResponse) LoadMore(ctx context.Context) (resp StreamingLocatorsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type StreamingLocatorsListOperationOptions struct {
	Filter  *string
	Orderby *string
	Top     *int64
}

func DefaultStreamingLocatorsListOperationOptions() StreamingLocatorsListOperationOptions {
	return StreamingLocatorsListOperationOptions{}
}

func (o StreamingLocatorsListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o StreamingLocatorsListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Orderby != nil {
		out["$orderby"] = *o.Orderby
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// StreamingLocatorsList ...
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingLocatorsList(ctx context.Context, id MediaServiceId, options StreamingLocatorsListOperationOptions) (resp StreamingLocatorsListOperationResponse, err error) {
	req, err := c.preparerForStreamingLocatorsList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForStreamingLocatorsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForStreamingLocatorsList prepares the StreamingLocatorsList request.
func (c StreamingPoliciesAndStreamingLocatorsClient) preparerForStreamingLocatorsList(ctx context.Context, id MediaServiceId, options StreamingLocatorsListOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/streamingLocators", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForStreamingLocatorsListWithNextLink prepares the StreamingLocatorsList request with the given nextLink token.
func (c StreamingPoliciesAndStreamingLocatorsClient) preparerForStreamingLocatorsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForStreamingLocatorsList handles the response to the StreamingLocatorsList request. The method always
// closes the http.Response Body.
func (c StreamingPoliciesAndStreamingLocatorsClient) responderForStreamingLocatorsList(resp *http.Response) (result StreamingLocatorsListOperationResponse, err error) {
	type page struct {
		Values   []StreamingLocator `json:"value"`
		NextLink *string            `json:"@odata.nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result StreamingLocatorsListOperationResponse, err error) {
			req, err := c.preparerForStreamingLocatorsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForStreamingLocatorsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "streamingpoliciesandstreaminglocators.StreamingPoliciesAndStreamingLocatorsClient", "StreamingLocatorsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// StreamingLocatorsListComplete retrieves all of the results into a single object
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingLocatorsListComplete(ctx context.Context, id MediaServiceId, options StreamingLocatorsListOperationOptions) (StreamingLocatorsListCompleteResult, error) {
	return c.StreamingLocatorsListCompleteMatchingPredicate(ctx, id, options, StreamingLocatorOperationPredicate{})
}

// StreamingLocatorsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c StreamingPoliciesAndStreamingLocatorsClient) StreamingLocatorsListCompleteMatchingPredicate(ctx context.Context, id MediaServiceId, options StreamingLocatorsListOperationOptions, predicate StreamingLocatorOperationPredicate) (resp StreamingLocatorsListCompleteResult, err error) {
	items := make([]StreamingLocator, 0)

	page, err := c.StreamingLocatorsList(ctx, id, options)
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

	out := StreamingLocatorsListCompleteResult{
		Items: items,
	}
	return out, nil
}
