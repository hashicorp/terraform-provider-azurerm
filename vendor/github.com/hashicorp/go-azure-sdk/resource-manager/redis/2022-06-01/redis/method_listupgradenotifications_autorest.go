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

type ListUpgradeNotificationsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]UpgradeNotification

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListUpgradeNotificationsOperationResponse, error)
}

type ListUpgradeNotificationsCompleteResult struct {
	Items []UpgradeNotification
}

func (r ListUpgradeNotificationsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListUpgradeNotificationsOperationResponse) LoadMore(ctx context.Context) (resp ListUpgradeNotificationsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListUpgradeNotificationsOperationOptions struct {
	History *float64
}

func DefaultListUpgradeNotificationsOperationOptions() ListUpgradeNotificationsOperationOptions {
	return ListUpgradeNotificationsOperationOptions{}
}

func (o ListUpgradeNotificationsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListUpgradeNotificationsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.History != nil {
		out["history"] = *o.History
	}

	return out
}

// ListUpgradeNotifications ...
func (c RedisClient) ListUpgradeNotifications(ctx context.Context, id RediId, options ListUpgradeNotificationsOperationOptions) (resp ListUpgradeNotificationsOperationResponse, err error) {
	req, err := c.preparerForListUpgradeNotifications(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "ListUpgradeNotifications", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "ListUpgradeNotifications", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListUpgradeNotifications(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "ListUpgradeNotifications", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListUpgradeNotifications prepares the ListUpgradeNotifications request.
func (c RedisClient) preparerForListUpgradeNotifications(ctx context.Context, id RediId, options ListUpgradeNotificationsOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/listUpgradeNotifications", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListUpgradeNotificationsWithNextLink prepares the ListUpgradeNotifications request with the given nextLink token.
func (c RedisClient) preparerForListUpgradeNotificationsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListUpgradeNotifications handles the response to the ListUpgradeNotifications request. The method always
// closes the http.Response Body.
func (c RedisClient) responderForListUpgradeNotifications(resp *http.Response) (result ListUpgradeNotificationsOperationResponse, err error) {
	type page struct {
		Values   []UpgradeNotification `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListUpgradeNotificationsOperationResponse, err error) {
			req, err := c.preparerForListUpgradeNotificationsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "redis.RedisClient", "ListUpgradeNotifications", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "redis.RedisClient", "ListUpgradeNotifications", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListUpgradeNotifications(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "redis.RedisClient", "ListUpgradeNotifications", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListUpgradeNotificationsComplete retrieves all of the results into a single object
func (c RedisClient) ListUpgradeNotificationsComplete(ctx context.Context, id RediId, options ListUpgradeNotificationsOperationOptions) (ListUpgradeNotificationsCompleteResult, error) {
	return c.ListUpgradeNotificationsCompleteMatchingPredicate(ctx, id, options, UpgradeNotificationOperationPredicate{})
}

// ListUpgradeNotificationsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c RedisClient) ListUpgradeNotificationsCompleteMatchingPredicate(ctx context.Context, id RediId, options ListUpgradeNotificationsOperationOptions, predicate UpgradeNotificationOperationPredicate) (resp ListUpgradeNotificationsCompleteResult, err error) {
	items := make([]UpgradeNotification, 0)

	page, err := c.ListUpgradeNotifications(ctx, id, options)
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

	out := ListUpgradeNotificationsCompleteResult{
		Items: items,
	}
	return out, nil
}
