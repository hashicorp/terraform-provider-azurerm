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

type ListByAccountOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Share

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByAccountOperationResponse, error)
}

type ListByAccountCompleteResult struct {
	Items []Share
}

func (r ListByAccountOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByAccountOperationResponse) LoadMore(ctx context.Context) (resp ListByAccountOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByAccountOperationOptions struct {
	Filter  *string
	Orderby *string
}

func DefaultListByAccountOperationOptions() ListByAccountOperationOptions {
	return ListByAccountOperationOptions{}
}

func (o ListByAccountOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListByAccountOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Orderby != nil {
		out["$orderby"] = *o.Orderby
	}

	return out
}

// ListByAccount ...
func (c ShareClient) ListByAccount(ctx context.Context, id AccountId, options ListByAccountOperationOptions) (resp ListByAccountOperationResponse, err error) {
	req, err := c.preparerForListByAccount(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ListByAccount", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ListByAccount", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByAccount(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ListByAccount", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByAccount prepares the ListByAccount request.
func (c ShareClient) preparerForListByAccount(ctx context.Context, id AccountId, options ListByAccountOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/shares", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByAccountWithNextLink prepares the ListByAccount request with the given nextLink token.
func (c ShareClient) preparerForListByAccountWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByAccount handles the response to the ListByAccount request. The method always
// closes the http.Response Body.
func (c ShareClient) responderForListByAccount(resp *http.Response) (result ListByAccountOperationResponse, err error) {
	type page struct {
		Values   []Share `json:"value"`
		NextLink *string `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByAccountOperationResponse, err error) {
			req, err := c.preparerForListByAccountWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "share.ShareClient", "ListByAccount", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "share.ShareClient", "ListByAccount", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByAccount(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "share.ShareClient", "ListByAccount", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByAccountComplete retrieves all of the results into a single object
func (c ShareClient) ListByAccountComplete(ctx context.Context, id AccountId, options ListByAccountOperationOptions) (ListByAccountCompleteResult, error) {
	return c.ListByAccountCompleteMatchingPredicate(ctx, id, options, ShareOperationPredicate{})
}

// ListByAccountCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ShareClient) ListByAccountCompleteMatchingPredicate(ctx context.Context, id AccountId, options ListByAccountOperationOptions, predicate ShareOperationPredicate) (resp ListByAccountCompleteResult, err error) {
	items := make([]Share, 0)

	page, err := c.ListByAccount(ctx, id, options)
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

	out := ListByAccountCompleteResult{
		Items: items,
	}
	return out, nil
}
