package rules

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

type ListBySubscriptionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Rule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListBySubscriptionsOperationResponse, error)
}

type ListBySubscriptionsCompleteResult struct {
	Items []Rule
}

func (r ListBySubscriptionsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListBySubscriptionsOperationResponse) LoadMore(ctx context.Context) (resp ListBySubscriptionsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListBySubscriptionsOperationOptions struct {
	Skip *int64
	Top  *int64
}

func DefaultListBySubscriptionsOperationOptions() ListBySubscriptionsOperationOptions {
	return ListBySubscriptionsOperationOptions{}
}

func (o ListBySubscriptionsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListBySubscriptionsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Skip != nil {
		out["$skip"] = *o.Skip
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListBySubscriptions ...
func (c RulesClient) ListBySubscriptions(ctx context.Context, id Subscriptions2Id, options ListBySubscriptionsOperationOptions) (resp ListBySubscriptionsOperationResponse, err error) {
	req, err := c.preparerForListBySubscriptions(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "ListBySubscriptions", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "ListBySubscriptions", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListBySubscriptions(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "rules.RulesClient", "ListBySubscriptions", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListBySubscriptionsComplete retrieves all of the results into a single object
func (c RulesClient) ListBySubscriptionsComplete(ctx context.Context, id Subscriptions2Id, options ListBySubscriptionsOperationOptions) (ListBySubscriptionsCompleteResult, error) {
	return c.ListBySubscriptionsCompleteMatchingPredicate(ctx, id, options, RuleOperationPredicate{})
}

// ListBySubscriptionsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c RulesClient) ListBySubscriptionsCompleteMatchingPredicate(ctx context.Context, id Subscriptions2Id, options ListBySubscriptionsOperationOptions, predicate RuleOperationPredicate) (resp ListBySubscriptionsCompleteResult, err error) {
	items := make([]Rule, 0)

	page, err := c.ListBySubscriptions(ctx, id, options)
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

	out := ListBySubscriptionsCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListBySubscriptions prepares the ListBySubscriptions request.
func (c RulesClient) preparerForListBySubscriptions(ctx context.Context, id Subscriptions2Id, options ListBySubscriptionsOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/rules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListBySubscriptionsWithNextLink prepares the ListBySubscriptions request with the given nextLink token.
func (c RulesClient) preparerForListBySubscriptionsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListBySubscriptions handles the response to the ListBySubscriptions request. The method always
// closes the http.Response Body.
func (c RulesClient) responderForListBySubscriptions(resp *http.Response) (result ListBySubscriptionsOperationResponse, err error) {
	type page struct {
		Values   []Rule  `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListBySubscriptionsOperationResponse, err error) {
			req, err := c.preparerForListBySubscriptionsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "rules.RulesClient", "ListBySubscriptions", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "rules.RulesClient", "ListBySubscriptions", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListBySubscriptions(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "rules.RulesClient", "ListBySubscriptions", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
