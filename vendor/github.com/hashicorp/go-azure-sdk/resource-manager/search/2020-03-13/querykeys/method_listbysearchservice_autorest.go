package querykeys

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

type ListBySearchServiceOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]QueryKey

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListBySearchServiceOperationResponse, error)
}

type ListBySearchServiceCompleteResult struct {
	Items []QueryKey
}

func (r ListBySearchServiceOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListBySearchServiceOperationResponse) LoadMore(ctx context.Context) (resp ListBySearchServiceOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListBySearchServiceOperationOptions struct {
	XMsClientRequestId *string
}

func DefaultListBySearchServiceOperationOptions() ListBySearchServiceOperationOptions {
	return ListBySearchServiceOperationOptions{}
}

func (o ListBySearchServiceOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.XMsClientRequestId != nil {
		out["x-ms-client-request-id"] = *o.XMsClientRequestId
	}

	return out
}

func (o ListBySearchServiceOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// ListBySearchService ...
func (c QueryKeysClient) ListBySearchService(ctx context.Context, id SearchServiceId, options ListBySearchServiceOperationOptions) (resp ListBySearchServiceOperationResponse, err error) {
	req, err := c.preparerForListBySearchService(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "querykeys.QueryKeysClient", "ListBySearchService", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "querykeys.QueryKeysClient", "ListBySearchService", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListBySearchService(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "querykeys.QueryKeysClient", "ListBySearchService", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListBySearchServiceComplete retrieves all of the results into a single object
func (c QueryKeysClient) ListBySearchServiceComplete(ctx context.Context, id SearchServiceId, options ListBySearchServiceOperationOptions) (ListBySearchServiceCompleteResult, error) {
	return c.ListBySearchServiceCompleteMatchingPredicate(ctx, id, options, QueryKeyOperationPredicate{})
}

// ListBySearchServiceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c QueryKeysClient) ListBySearchServiceCompleteMatchingPredicate(ctx context.Context, id SearchServiceId, options ListBySearchServiceOperationOptions, predicate QueryKeyOperationPredicate) (resp ListBySearchServiceCompleteResult, err error) {
	items := make([]QueryKey, 0)

	page, err := c.ListBySearchService(ctx, id, options)
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

	out := ListBySearchServiceCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListBySearchService prepares the ListBySearchService request.
func (c QueryKeysClient) preparerForListBySearchService(ctx context.Context, id SearchServiceId, options ListBySearchServiceOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/listQueryKeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListBySearchServiceWithNextLink prepares the ListBySearchService request with the given nextLink token.
func (c QueryKeysClient) preparerForListBySearchServiceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListBySearchService handles the response to the ListBySearchService request. The method always
// closes the http.Response Body.
func (c QueryKeysClient) responderForListBySearchService(resp *http.Response) (result ListBySearchServiceOperationResponse, err error) {
	type page struct {
		Values   []QueryKey `json:"value"`
		NextLink *string    `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListBySearchServiceOperationResponse, err error) {
			req, err := c.preparerForListBySearchServiceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "querykeys.QueryKeysClient", "ListBySearchService", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "querykeys.QueryKeysClient", "ListBySearchService", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListBySearchService(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "querykeys.QueryKeysClient", "ListBySearchService", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
