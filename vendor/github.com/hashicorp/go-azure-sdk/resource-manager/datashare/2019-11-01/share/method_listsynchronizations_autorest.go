package share

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

type ListSynchronizationsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ShareSynchronization

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListSynchronizationsOperationResponse, error)
}

type ListSynchronizationsCompleteResult struct {
	Items []ShareSynchronization
}

func (r ListSynchronizationsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListSynchronizationsOperationResponse) LoadMore(ctx context.Context) (resp ListSynchronizationsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListSynchronizationsOperationOptions struct {
	Filter  *string
	Orderby *string
}

func DefaultListSynchronizationsOperationOptions() ListSynchronizationsOperationOptions {
	return ListSynchronizationsOperationOptions{}
}

func (o ListSynchronizationsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListSynchronizationsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Orderby != nil {
		out["$orderby"] = *o.Orderby
	}

	return out
}

// ListSynchronizations ...
func (c ShareClient) ListSynchronizations(ctx context.Context, id ShareId, options ListSynchronizationsOperationOptions) (resp ListSynchronizationsOperationResponse, err error) {
	req, err := c.preparerForListSynchronizations(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ListSynchronizations", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ListSynchronizations", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListSynchronizations(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ListSynchronizations", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListSynchronizations prepares the ListSynchronizations request.
func (c ShareClient) preparerForListSynchronizations(ctx context.Context, id ShareId, options ListSynchronizationsOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/listSynchronizations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListSynchronizationsWithNextLink prepares the ListSynchronizations request with the given nextLink token.
func (c ShareClient) preparerForListSynchronizationsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListSynchronizations handles the response to the ListSynchronizations request. The method always
// closes the http.Response Body.
func (c ShareClient) responderForListSynchronizations(resp *http.Response) (result ListSynchronizationsOperationResponse, err error) {
	type page struct {
		Values   []ShareSynchronization `json:"value"`
		NextLink *string                `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListSynchronizationsOperationResponse, err error) {
			req, err := c.preparerForListSynchronizationsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "share.ShareClient", "ListSynchronizations", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "share.ShareClient", "ListSynchronizations", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListSynchronizations(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "share.ShareClient", "ListSynchronizations", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListSynchronizationsComplete retrieves all of the results into a single object
func (c ShareClient) ListSynchronizationsComplete(ctx context.Context, id ShareId, options ListSynchronizationsOperationOptions) (ListSynchronizationsCompleteResult, error) {
	return c.ListSynchronizationsCompleteMatchingPredicate(ctx, id, options, ShareSynchronizationOperationPredicate{})
}

// ListSynchronizationsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ShareClient) ListSynchronizationsCompleteMatchingPredicate(ctx context.Context, id ShareId, options ListSynchronizationsOperationOptions, predicate ShareSynchronizationOperationPredicate) (resp ListSynchronizationsCompleteResult, err error) {
	items := make([]ShareSynchronization, 0)

	page, err := c.ListSynchronizations(ctx, id, options)
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

	out := ListSynchronizationsCompleteResult{
		Items: items,
	}
	return out, nil
}
