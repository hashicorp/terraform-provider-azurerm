package redis

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

type FirewallRulesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]RedisFirewallRule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (FirewallRulesListOperationResponse, error)
}

type FirewallRulesListCompleteResult struct {
	Items []RedisFirewallRule
}

func (r FirewallRulesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r FirewallRulesListOperationResponse) LoadMore(ctx context.Context) (resp FirewallRulesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// FirewallRulesList ...
func (c RedisClient) FirewallRulesList(ctx context.Context, id RediId) (resp FirewallRulesListOperationResponse, err error) {
	req, err := c.preparerForFirewallRulesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "FirewallRulesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "FirewallRulesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForFirewallRulesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "FirewallRulesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForFirewallRulesList prepares the FirewallRulesList request.
func (c RedisClient) preparerForFirewallRulesList(ctx context.Context, id RediId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/firewallRules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForFirewallRulesListWithNextLink prepares the FirewallRulesList request with the given nextLink token.
func (c RedisClient) preparerForFirewallRulesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForFirewallRulesList handles the response to the FirewallRulesList request. The method always
// closes the http.Response Body.
func (c RedisClient) responderForFirewallRulesList(resp *http.Response) (result FirewallRulesListOperationResponse, err error) {
	type page struct {
		Values   []RedisFirewallRule `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result FirewallRulesListOperationResponse, err error) {
			req, err := c.preparerForFirewallRulesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "redis.RedisClient", "FirewallRulesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "redis.RedisClient", "FirewallRulesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForFirewallRulesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "redis.RedisClient", "FirewallRulesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// FirewallRulesListComplete retrieves all of the results into a single object
func (c RedisClient) FirewallRulesListComplete(ctx context.Context, id RediId) (FirewallRulesListCompleteResult, error) {
	return c.FirewallRulesListCompleteMatchingPredicate(ctx, id, RedisFirewallRuleOperationPredicate{})
}

// FirewallRulesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c RedisClient) FirewallRulesListCompleteMatchingPredicate(ctx context.Context, id RediId, predicate RedisFirewallRuleOperationPredicate) (resp FirewallRulesListCompleteResult, err error) {
	items := make([]RedisFirewallRule, 0)

	page, err := c.FirewallRulesList(ctx, id)
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

	out := FirewallRulesListCompleteResult{
		Items: items,
	}
	return out, nil
}
