package subscriptions

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

type ListByTopicOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SBSubscription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByTopicOperationResponse, error)
}

type ListByTopicCompleteResult struct {
	Items []SBSubscription
}

func (r ListByTopicOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByTopicOperationResponse) LoadMore(ctx context.Context) (resp ListByTopicOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByTopicOperationOptions struct {
	Skip *int64
	Top  *int64
}

func DefaultListByTopicOperationOptions() ListByTopicOperationOptions {
	return ListByTopicOperationOptions{}
}

func (o ListByTopicOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListByTopicOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Skip != nil {
		out["$skip"] = *o.Skip
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListByTopic ...
func (c SubscriptionsClient) ListByTopic(ctx context.Context, id TopicId, options ListByTopicOperationOptions) (resp ListByTopicOperationResponse, err error) {
	req, err := c.preparerForListByTopic(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "ListByTopic", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "ListByTopic", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByTopic(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "ListByTopic", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByTopicComplete retrieves all of the results into a single object
func (c SubscriptionsClient) ListByTopicComplete(ctx context.Context, id TopicId, options ListByTopicOperationOptions) (ListByTopicCompleteResult, error) {
	return c.ListByTopicCompleteMatchingPredicate(ctx, id, options, SBSubscriptionOperationPredicate{})
}

// ListByTopicCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SubscriptionsClient) ListByTopicCompleteMatchingPredicate(ctx context.Context, id TopicId, options ListByTopicOperationOptions, predicate SBSubscriptionOperationPredicate) (resp ListByTopicCompleteResult, err error) {
	items := make([]SBSubscription, 0)

	page, err := c.ListByTopic(ctx, id, options)
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

	out := ListByTopicCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByTopic prepares the ListByTopic request.
func (c SubscriptionsClient) preparerForListByTopic(ctx context.Context, id TopicId, options ListByTopicOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/subscriptions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByTopicWithNextLink prepares the ListByTopic request with the given nextLink token.
func (c SubscriptionsClient) preparerForListByTopicWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByTopic handles the response to the ListByTopic request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForListByTopic(resp *http.Response) (result ListByTopicOperationResponse, err error) {
	type page struct {
		Values   []SBSubscription `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByTopicOperationResponse, err error) {
			req, err := c.preparerForListByTopicWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "ListByTopic", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "ListByTopic", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByTopic(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "ListByTopic", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
