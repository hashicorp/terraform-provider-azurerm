package storagetargets

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

type ListByCacheOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]StorageTarget

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByCacheOperationResponse, error)
}

type ListByCacheCompleteResult struct {
	Items []StorageTarget
}

func (r ListByCacheOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByCacheOperationResponse) LoadMore(ctx context.Context) (resp ListByCacheOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByCache ...
func (c StorageTargetsClient) ListByCache(ctx context.Context, id CacheId) (resp ListByCacheOperationResponse, err error) {
	req, err := c.preparerForListByCache(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "ListByCache", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "ListByCache", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByCache(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "ListByCache", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByCache prepares the ListByCache request.
func (c StorageTargetsClient) preparerForListByCache(ctx context.Context, id CacheId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/storageTargets", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByCacheWithNextLink prepares the ListByCache request with the given nextLink token.
func (c StorageTargetsClient) preparerForListByCacheWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByCache handles the response to the ListByCache request. The method always
// closes the http.Response Body.
func (c StorageTargetsClient) responderForListByCache(resp *http.Response) (result ListByCacheOperationResponse, err error) {
	type page struct {
		Values   []StorageTarget `json:"value"`
		NextLink *string         `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByCacheOperationResponse, err error) {
			req, err := c.preparerForListByCacheWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "ListByCache", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "ListByCache", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByCache(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "storagetargets.StorageTargetsClient", "ListByCache", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByCacheComplete retrieves all of the results into a single object
func (c StorageTargetsClient) ListByCacheComplete(ctx context.Context, id CacheId) (ListByCacheCompleteResult, error) {
	return c.ListByCacheCompleteMatchingPredicate(ctx, id, StorageTargetOperationPredicate{})
}

// ListByCacheCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c StorageTargetsClient) ListByCacheCompleteMatchingPredicate(ctx context.Context, id CacheId, predicate StorageTargetOperationPredicate) (resp ListByCacheCompleteResult, err error) {
	items := make([]StorageTarget, 0)

	page, err := c.ListByCache(ctx, id)
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

	out := ListByCacheCompleteResult{
		Items: items,
	}
	return out, nil
}
