package views

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByScopeOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]View

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByScopeOperationResponse, error)
}

type ListByScopeCompleteResult struct {
	Items []View
}

func (r ListByScopeOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByScopeOperationResponse) LoadMore(ctx context.Context) (resp ListByScopeOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByScope ...
func (c ViewsClient) ListByScope(ctx context.Context, id commonids.ScopeId) (resp ListByScopeOperationResponse, err error) {
	req, err := c.preparerForListByScope(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "ListByScope", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "ListByScope", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByScope(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "views.ViewsClient", "ListByScope", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByScope prepares the ListByScope request.
func (c ViewsClient) preparerForListByScope(ctx context.Context, id commonids.ScopeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.CostManagement/views", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByScopeWithNextLink prepares the ListByScope request with the given nextLink token.
func (c ViewsClient) preparerForListByScopeWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByScope handles the response to the ListByScope request. The method always
// closes the http.Response Body.
func (c ViewsClient) responderForListByScope(resp *http.Response) (result ListByScopeOperationResponse, err error) {
	type page struct {
		Values   []View  `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByScopeOperationResponse, err error) {
			req, err := c.preparerForListByScopeWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "views.ViewsClient", "ListByScope", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "views.ViewsClient", "ListByScope", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByScope(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "views.ViewsClient", "ListByScope", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByScopeComplete retrieves all of the results into a single object
func (c ViewsClient) ListByScopeComplete(ctx context.Context, id commonids.ScopeId) (ListByScopeCompleteResult, error) {
	return c.ListByScopeCompleteMatchingPredicate(ctx, id, ViewOperationPredicate{})
}

// ListByScopeCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ViewsClient) ListByScopeCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, predicate ViewOperationPredicate) (resp ListByScopeCompleteResult, err error) {
	items := make([]View, 0)

	page, err := c.ListByScope(ctx, id)
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

	out := ListByScopeCompleteResult{
		Items: items,
	}
	return out, nil
}
