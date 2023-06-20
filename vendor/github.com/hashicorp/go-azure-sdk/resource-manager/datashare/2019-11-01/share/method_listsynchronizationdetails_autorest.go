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

type ListSynchronizationDetailsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SynchronizationDetails

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListSynchronizationDetailsOperationResponse, error)
}

type ListSynchronizationDetailsCompleteResult struct {
	Items []SynchronizationDetails
}

func (r ListSynchronizationDetailsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListSynchronizationDetailsOperationResponse) LoadMore(ctx context.Context) (resp ListSynchronizationDetailsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListSynchronizationDetailsOperationOptions struct {
	Filter  *string
	Orderby *string
}

func DefaultListSynchronizationDetailsOperationOptions() ListSynchronizationDetailsOperationOptions {
	return ListSynchronizationDetailsOperationOptions{}
}

func (o ListSynchronizationDetailsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListSynchronizationDetailsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Orderby != nil {
		out["$orderby"] = *o.Orderby
	}

	return out
}

// ListSynchronizationDetails ...
func (c ShareClient) ListSynchronizationDetails(ctx context.Context, id ShareId, input ShareSynchronization, options ListSynchronizationDetailsOperationOptions) (resp ListSynchronizationDetailsOperationResponse, err error) {
	req, err := c.preparerForListSynchronizationDetails(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ListSynchronizationDetails", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ListSynchronizationDetails", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListSynchronizationDetails(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ListSynchronizationDetails", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListSynchronizationDetails prepares the ListSynchronizationDetails request.
func (c ShareClient) preparerForListSynchronizationDetails(ctx context.Context, id ShareId, input ShareSynchronization, options ListSynchronizationDetailsOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/listSynchronizationDetails", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListSynchronizationDetailsWithNextLink prepares the ListSynchronizationDetails request with the given nextLink token.
func (c ShareClient) preparerForListSynchronizationDetailsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListSynchronizationDetails handles the response to the ListSynchronizationDetails request. The method always
// closes the http.Response Body.
func (c ShareClient) responderForListSynchronizationDetails(resp *http.Response) (result ListSynchronizationDetailsOperationResponse, err error) {
	type page struct {
		Values   []SynchronizationDetails `json:"value"`
		NextLink *string                  `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListSynchronizationDetailsOperationResponse, err error) {
			req, err := c.preparerForListSynchronizationDetailsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "share.ShareClient", "ListSynchronizationDetails", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "share.ShareClient", "ListSynchronizationDetails", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListSynchronizationDetails(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "share.ShareClient", "ListSynchronizationDetails", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListSynchronizationDetailsComplete retrieves all of the results into a single object
func (c ShareClient) ListSynchronizationDetailsComplete(ctx context.Context, id ShareId, input ShareSynchronization, options ListSynchronizationDetailsOperationOptions) (ListSynchronizationDetailsCompleteResult, error) {
	return c.ListSynchronizationDetailsCompleteMatchingPredicate(ctx, id, input, options, SynchronizationDetailsOperationPredicate{})
}

// ListSynchronizationDetailsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ShareClient) ListSynchronizationDetailsCompleteMatchingPredicate(ctx context.Context, id ShareId, input ShareSynchronization, options ListSynchronizationDetailsOperationOptions, predicate SynchronizationDetailsOperationPredicate) (resp ListSynchronizationDetailsCompleteResult, err error) {
	items := make([]SynchronizationDetails, 0)

	page, err := c.ListSynchronizationDetails(ctx, id, input, options)
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

	out := ListSynchronizationDetailsCompleteResult{
		Items: items,
	}
	return out, nil
}
