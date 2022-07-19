package namespaces

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

type ListNetworkRuleSetsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]NetworkRuleSet

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListNetworkRuleSetsOperationResponse, error)
}

type ListNetworkRuleSetsCompleteResult struct {
	Items []NetworkRuleSet
}

func (r ListNetworkRuleSetsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListNetworkRuleSetsOperationResponse) LoadMore(ctx context.Context) (resp ListNetworkRuleSetsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListNetworkRuleSets ...
func (c NamespacesClient) ListNetworkRuleSets(ctx context.Context, id NamespaceId) (resp ListNetworkRuleSetsOperationResponse, err error) {
	req, err := c.preparerForListNetworkRuleSets(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "ListNetworkRuleSets", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "ListNetworkRuleSets", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListNetworkRuleSets(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "ListNetworkRuleSets", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListNetworkRuleSetsComplete retrieves all of the results into a single object
func (c NamespacesClient) ListNetworkRuleSetsComplete(ctx context.Context, id NamespaceId) (ListNetworkRuleSetsCompleteResult, error) {
	return c.ListNetworkRuleSetsCompleteMatchingPredicate(ctx, id, NetworkRuleSetOperationPredicate{})
}

// ListNetworkRuleSetsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c NamespacesClient) ListNetworkRuleSetsCompleteMatchingPredicate(ctx context.Context, id NamespaceId, predicate NetworkRuleSetOperationPredicate) (resp ListNetworkRuleSetsCompleteResult, err error) {
	items := make([]NetworkRuleSet, 0)

	page, err := c.ListNetworkRuleSets(ctx, id)
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

	out := ListNetworkRuleSetsCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListNetworkRuleSets prepares the ListNetworkRuleSets request.
func (c NamespacesClient) preparerForListNetworkRuleSets(ctx context.Context, id NamespaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/networkRuleSets", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListNetworkRuleSetsWithNextLink prepares the ListNetworkRuleSets request with the given nextLink token.
func (c NamespacesClient) preparerForListNetworkRuleSetsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListNetworkRuleSets handles the response to the ListNetworkRuleSets request. The method always
// closes the http.Response Body.
func (c NamespacesClient) responderForListNetworkRuleSets(resp *http.Response) (result ListNetworkRuleSetsOperationResponse, err error) {
	type page struct {
		Values   []NetworkRuleSet `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListNetworkRuleSetsOperationResponse, err error) {
			req, err := c.preparerForListNetworkRuleSetsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "ListNetworkRuleSets", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "ListNetworkRuleSets", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListNetworkRuleSets(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "ListNetworkRuleSets", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
