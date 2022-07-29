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

type ListAuthorizationRulesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AuthorizationRule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListAuthorizationRulesOperationResponse, error)
}

type ListAuthorizationRulesCompleteResult struct {
	Items []AuthorizationRule
}

func (r ListAuthorizationRulesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListAuthorizationRulesOperationResponse) LoadMore(ctx context.Context) (resp ListAuthorizationRulesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListAuthorizationRules ...
func (c NamespacesClient) ListAuthorizationRules(ctx context.Context, id NamespaceId) (resp ListAuthorizationRulesOperationResponse, err error) {
	req, err := c.preparerForListAuthorizationRules(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "ListAuthorizationRules", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "ListAuthorizationRules", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListAuthorizationRules(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "ListAuthorizationRules", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListAuthorizationRulesComplete retrieves all of the results into a single object
func (c NamespacesClient) ListAuthorizationRulesComplete(ctx context.Context, id NamespaceId) (ListAuthorizationRulesCompleteResult, error) {
	return c.ListAuthorizationRulesCompleteMatchingPredicate(ctx, id, AuthorizationRuleOperationPredicate{})
}

// ListAuthorizationRulesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c NamespacesClient) ListAuthorizationRulesCompleteMatchingPredicate(ctx context.Context, id NamespaceId, predicate AuthorizationRuleOperationPredicate) (resp ListAuthorizationRulesCompleteResult, err error) {
	items := make([]AuthorizationRule, 0)

	page, err := c.ListAuthorizationRules(ctx, id)
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

	out := ListAuthorizationRulesCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListAuthorizationRules prepares the ListAuthorizationRules request.
func (c NamespacesClient) preparerForListAuthorizationRules(ctx context.Context, id NamespaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/authorizationRules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListAuthorizationRulesWithNextLink prepares the ListAuthorizationRules request with the given nextLink token.
func (c NamespacesClient) preparerForListAuthorizationRulesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListAuthorizationRules handles the response to the ListAuthorizationRules request. The method always
// closes the http.Response Body.
func (c NamespacesClient) responderForListAuthorizationRules(resp *http.Response) (result ListAuthorizationRulesOperationResponse, err error) {
	type page struct {
		Values   []AuthorizationRule `json:"value"`
		NextLink *string             `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListAuthorizationRulesOperationResponse, err error) {
			req, err := c.preparerForListAuthorizationRulesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "ListAuthorizationRules", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "ListAuthorizationRules", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListAuthorizationRules(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "namespaces.NamespacesClient", "ListAuthorizationRules", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
