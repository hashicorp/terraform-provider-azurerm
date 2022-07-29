package queuesauthorizationrule

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

type QueuesListAuthorizationRulesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SBAuthorizationRule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (QueuesListAuthorizationRulesOperationResponse, error)
}

type QueuesListAuthorizationRulesCompleteResult struct {
	Items []SBAuthorizationRule
}

func (r QueuesListAuthorizationRulesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r QueuesListAuthorizationRulesOperationResponse) LoadMore(ctx context.Context) (resp QueuesListAuthorizationRulesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// QueuesListAuthorizationRules ...
func (c QueuesAuthorizationRuleClient) QueuesListAuthorizationRules(ctx context.Context, id QueueId) (resp QueuesListAuthorizationRulesOperationResponse, err error) {
	req, err := c.preparerForQueuesListAuthorizationRules(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesListAuthorizationRules", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesListAuthorizationRules", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForQueuesListAuthorizationRules(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesListAuthorizationRules", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// QueuesListAuthorizationRulesComplete retrieves all of the results into a single object
func (c QueuesAuthorizationRuleClient) QueuesListAuthorizationRulesComplete(ctx context.Context, id QueueId) (QueuesListAuthorizationRulesCompleteResult, error) {
	return c.QueuesListAuthorizationRulesCompleteMatchingPredicate(ctx, id, SBAuthorizationRuleOperationPredicate{})
}

// QueuesListAuthorizationRulesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c QueuesAuthorizationRuleClient) QueuesListAuthorizationRulesCompleteMatchingPredicate(ctx context.Context, id QueueId, predicate SBAuthorizationRuleOperationPredicate) (resp QueuesListAuthorizationRulesCompleteResult, err error) {
	items := make([]SBAuthorizationRule, 0)

	page, err := c.QueuesListAuthorizationRules(ctx, id)
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

	out := QueuesListAuthorizationRulesCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForQueuesListAuthorizationRules prepares the QueuesListAuthorizationRules request.
func (c QueuesAuthorizationRuleClient) preparerForQueuesListAuthorizationRules(ctx context.Context, id QueueId) (*http.Request, error) {
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

// preparerForQueuesListAuthorizationRulesWithNextLink prepares the QueuesListAuthorizationRules request with the given nextLink token.
func (c QueuesAuthorizationRuleClient) preparerForQueuesListAuthorizationRulesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForQueuesListAuthorizationRules handles the response to the QueuesListAuthorizationRules request. The method always
// closes the http.Response Body.
func (c QueuesAuthorizationRuleClient) responderForQueuesListAuthorizationRules(resp *http.Response) (result QueuesListAuthorizationRulesOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result QueuesListAuthorizationRulesOperationResponse, err error) {
			req, err := c.preparerForQueuesListAuthorizationRulesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesListAuthorizationRules", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesListAuthorizationRules", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForQueuesListAuthorizationRules(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "queuesauthorizationrule.QueuesAuthorizationRuleClient", "QueuesListAuthorizationRules", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
