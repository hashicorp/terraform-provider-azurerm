package namespacesauthorizationrule

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

type NamespacesListAuthorizationRulesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SBAuthorizationRule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (NamespacesListAuthorizationRulesOperationResponse, error)
}

type NamespacesListAuthorizationRulesCompleteResult struct {
	Items []SBAuthorizationRule
}

func (r NamespacesListAuthorizationRulesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r NamespacesListAuthorizationRulesOperationResponse) LoadMore(ctx context.Context) (resp NamespacesListAuthorizationRulesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// NamespacesListAuthorizationRules ...
func (c NamespacesAuthorizationRuleClient) NamespacesListAuthorizationRules(ctx context.Context, id NamespaceId) (resp NamespacesListAuthorizationRulesOperationResponse, err error) {
	req, err := c.preparerForNamespacesListAuthorizationRules(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespacesauthorizationrule.NamespacesAuthorizationRuleClient", "NamespacesListAuthorizationRules", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespacesauthorizationrule.NamespacesAuthorizationRuleClient", "NamespacesListAuthorizationRules", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForNamespacesListAuthorizationRules(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "namespacesauthorizationrule.NamespacesAuthorizationRuleClient", "NamespacesListAuthorizationRules", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// NamespacesListAuthorizationRulesComplete retrieves all of the results into a single object
func (c NamespacesAuthorizationRuleClient) NamespacesListAuthorizationRulesComplete(ctx context.Context, id NamespaceId) (NamespacesListAuthorizationRulesCompleteResult, error) {
	return c.NamespacesListAuthorizationRulesCompleteMatchingPredicate(ctx, id, SBAuthorizationRuleOperationPredicate{})
}

// NamespacesListAuthorizationRulesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c NamespacesAuthorizationRuleClient) NamespacesListAuthorizationRulesCompleteMatchingPredicate(ctx context.Context, id NamespaceId, predicate SBAuthorizationRuleOperationPredicate) (resp NamespacesListAuthorizationRulesCompleteResult, err error) {
	items := make([]SBAuthorizationRule, 0)

	page, err := c.NamespacesListAuthorizationRules(ctx, id)
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

	out := NamespacesListAuthorizationRulesCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForNamespacesListAuthorizationRules prepares the NamespacesListAuthorizationRules request.
func (c NamespacesAuthorizationRuleClient) preparerForNamespacesListAuthorizationRules(ctx context.Context, id NamespaceId) (*http.Request, error) {
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

// preparerForNamespacesListAuthorizationRulesWithNextLink prepares the NamespacesListAuthorizationRules request with the given nextLink token.
func (c NamespacesAuthorizationRuleClient) preparerForNamespacesListAuthorizationRulesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForNamespacesListAuthorizationRules handles the response to the NamespacesListAuthorizationRules request. The method always
// closes the http.Response Body.
func (c NamespacesAuthorizationRuleClient) responderForNamespacesListAuthorizationRules(resp *http.Response) (result NamespacesListAuthorizationRulesOperationResponse, err error) {
	type page struct {
		Values   []SBAuthorizationRule `json:"value"`
		NextLink *string               `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result NamespacesListAuthorizationRulesOperationResponse, err error) {
			req, err := c.preparerForNamespacesListAuthorizationRulesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "namespacesauthorizationrule.NamespacesAuthorizationRuleClient", "NamespacesListAuthorizationRules", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "namespacesauthorizationrule.NamespacesAuthorizationRuleClient", "NamespacesListAuthorizationRules", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForNamespacesListAuthorizationRules(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "namespacesauthorizationrule.NamespacesAuthorizationRuleClient", "NamespacesListAuthorizationRules", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
