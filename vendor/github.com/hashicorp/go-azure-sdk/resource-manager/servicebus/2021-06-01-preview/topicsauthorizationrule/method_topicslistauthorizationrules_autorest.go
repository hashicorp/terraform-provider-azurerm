package topicsauthorizationrule

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

type TopicsListAuthorizationRulesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SBAuthorizationRule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (TopicsListAuthorizationRulesOperationResponse, error)
}

type TopicsListAuthorizationRulesCompleteResult struct {
	Items []SBAuthorizationRule
}

func (r TopicsListAuthorizationRulesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r TopicsListAuthorizationRulesOperationResponse) LoadMore(ctx context.Context) (resp TopicsListAuthorizationRulesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// TopicsListAuthorizationRules ...
func (c TopicsAuthorizationRuleClient) TopicsListAuthorizationRules(ctx context.Context, id TopicId) (resp TopicsListAuthorizationRulesOperationResponse, err error) {
	req, err := c.preparerForTopicsListAuthorizationRules(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsListAuthorizationRules", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsListAuthorizationRules", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForTopicsListAuthorizationRules(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsListAuthorizationRules", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// TopicsListAuthorizationRulesComplete retrieves all of the results into a single object
func (c TopicsAuthorizationRuleClient) TopicsListAuthorizationRulesComplete(ctx context.Context, id TopicId) (TopicsListAuthorizationRulesCompleteResult, error) {
	return c.TopicsListAuthorizationRulesCompleteMatchingPredicate(ctx, id, SBAuthorizationRuleOperationPredicate{})
}

// TopicsListAuthorizationRulesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c TopicsAuthorizationRuleClient) TopicsListAuthorizationRulesCompleteMatchingPredicate(ctx context.Context, id TopicId, predicate SBAuthorizationRuleOperationPredicate) (resp TopicsListAuthorizationRulesCompleteResult, err error) {
	items := make([]SBAuthorizationRule, 0)

	page, err := c.TopicsListAuthorizationRules(ctx, id)
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

	out := TopicsListAuthorizationRulesCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForTopicsListAuthorizationRules prepares the TopicsListAuthorizationRules request.
func (c TopicsAuthorizationRuleClient) preparerForTopicsListAuthorizationRules(ctx context.Context, id TopicId) (*http.Request, error) {
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

// preparerForTopicsListAuthorizationRulesWithNextLink prepares the TopicsListAuthorizationRules request with the given nextLink token.
func (c TopicsAuthorizationRuleClient) preparerForTopicsListAuthorizationRulesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForTopicsListAuthorizationRules handles the response to the TopicsListAuthorizationRules request. The method always
// closes the http.Response Body.
func (c TopicsAuthorizationRuleClient) responderForTopicsListAuthorizationRules(resp *http.Response) (result TopicsListAuthorizationRulesOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result TopicsListAuthorizationRulesOperationResponse, err error) {
			req, err := c.preparerForTopicsListAuthorizationRulesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsListAuthorizationRules", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsListAuthorizationRules", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForTopicsListAuthorizationRules(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "topicsauthorizationrule.TopicsAuthorizationRuleClient", "TopicsListAuthorizationRules", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
