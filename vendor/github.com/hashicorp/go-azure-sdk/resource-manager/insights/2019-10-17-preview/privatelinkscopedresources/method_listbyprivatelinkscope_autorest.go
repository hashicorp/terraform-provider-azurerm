package privatelinkscopedresources

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

type ListByPrivateLinkScopeOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ScopedResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByPrivateLinkScopeOperationResponse, error)
}

type ListByPrivateLinkScopeCompleteResult struct {
	Items []ScopedResource
}

func (r ListByPrivateLinkScopeOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByPrivateLinkScopeOperationResponse) LoadMore(ctx context.Context) (resp ListByPrivateLinkScopeOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByPrivateLinkScope ...
func (c PrivateLinkScopedResourcesClient) ListByPrivateLinkScope(ctx context.Context, id PrivateLinkScopeId) (resp ListByPrivateLinkScopeOperationResponse, err error) {
	req, err := c.preparerForListByPrivateLinkScope(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopedresources.PrivateLinkScopedResourcesClient", "ListByPrivateLinkScope", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopedresources.PrivateLinkScopedResourcesClient", "ListByPrivateLinkScope", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByPrivateLinkScope(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopedresources.PrivateLinkScopedResourcesClient", "ListByPrivateLinkScope", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByPrivateLinkScope prepares the ListByPrivateLinkScope request.
func (c PrivateLinkScopedResourcesClient) preparerForListByPrivateLinkScope(ctx context.Context, id PrivateLinkScopeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/scopedResources", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByPrivateLinkScopeWithNextLink prepares the ListByPrivateLinkScope request with the given nextLink token.
func (c PrivateLinkScopedResourcesClient) preparerForListByPrivateLinkScopeWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByPrivateLinkScope handles the response to the ListByPrivateLinkScope request. The method always
// closes the http.Response Body.
func (c PrivateLinkScopedResourcesClient) responderForListByPrivateLinkScope(resp *http.Response) (result ListByPrivateLinkScopeOperationResponse, err error) {
	type page struct {
		Values   []ScopedResource `json:"value"`
		NextLink *string          `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByPrivateLinkScopeOperationResponse, err error) {
			req, err := c.preparerForListByPrivateLinkScopeWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "privatelinkscopedresources.PrivateLinkScopedResourcesClient", "ListByPrivateLinkScope", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "privatelinkscopedresources.PrivateLinkScopedResourcesClient", "ListByPrivateLinkScope", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByPrivateLinkScope(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "privatelinkscopedresources.PrivateLinkScopedResourcesClient", "ListByPrivateLinkScope", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByPrivateLinkScopeComplete retrieves all of the results into a single object
func (c PrivateLinkScopedResourcesClient) ListByPrivateLinkScopeComplete(ctx context.Context, id PrivateLinkScopeId) (ListByPrivateLinkScopeCompleteResult, error) {
	return c.ListByPrivateLinkScopeCompleteMatchingPredicate(ctx, id, ScopedResourceOperationPredicate{})
}

// ListByPrivateLinkScopeCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c PrivateLinkScopedResourcesClient) ListByPrivateLinkScopeCompleteMatchingPredicate(ctx context.Context, id PrivateLinkScopeId, predicate ScopedResourceOperationPredicate) (resp ListByPrivateLinkScopeCompleteResult, err error) {
	items := make([]ScopedResource, 0)

	page, err := c.ListByPrivateLinkScope(ctx, id)
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

	out := ListByPrivateLinkScopeCompleteResult{
		Items: items,
	}
	return out, nil
}
