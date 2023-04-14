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

type LinkedServerListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]RedisLinkedServerWithProperties

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (LinkedServerListOperationResponse, error)
}

type LinkedServerListCompleteResult struct {
	Items []RedisLinkedServerWithProperties
}

func (r LinkedServerListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r LinkedServerListOperationResponse) LoadMore(ctx context.Context) (resp LinkedServerListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// LinkedServerList ...
func (c RedisClient) LinkedServerList(ctx context.Context, id RediId) (resp LinkedServerListOperationResponse, err error) {
	req, err := c.preparerForLinkedServerList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForLinkedServerList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForLinkedServerList prepares the LinkedServerList request.
func (c RedisClient) preparerForLinkedServerList(ctx context.Context, id RediId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/linkedServers", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForLinkedServerListWithNextLink prepares the LinkedServerList request with the given nextLink token.
func (c RedisClient) preparerForLinkedServerListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForLinkedServerList handles the response to the LinkedServerList request. The method always
// closes the http.Response Body.
func (c RedisClient) responderForLinkedServerList(resp *http.Response) (result LinkedServerListOperationResponse, err error) {
	type page struct {
		Values   []RedisLinkedServerWithProperties `json:"value"`
		NextLink *string                           `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result LinkedServerListOperationResponse, err error) {
			req, err := c.preparerForLinkedServerListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForLinkedServerList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "redis.RedisClient", "LinkedServerList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// LinkedServerListComplete retrieves all of the results into a single object
func (c RedisClient) LinkedServerListComplete(ctx context.Context, id RediId) (LinkedServerListCompleteResult, error) {
	return c.LinkedServerListCompleteMatchingPredicate(ctx, id, RedisLinkedServerWithPropertiesOperationPredicate{})
}

// LinkedServerListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c RedisClient) LinkedServerListCompleteMatchingPredicate(ctx context.Context, id RediId, predicate RedisLinkedServerWithPropertiesOperationPredicate) (resp LinkedServerListCompleteResult, err error) {
	items := make([]RedisLinkedServerWithProperties, 0)

	page, err := c.LinkedServerList(ctx, id)
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

	out := LinkedServerListCompleteResult{
		Items: items,
	}
	return out, nil
}
