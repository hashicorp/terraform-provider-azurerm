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

type PatchSchedulesListByRedisResourceOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]RedisPatchSchedule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (PatchSchedulesListByRedisResourceOperationResponse, error)
}

type PatchSchedulesListByRedisResourceCompleteResult struct {
	Items []RedisPatchSchedule
}

func (r PatchSchedulesListByRedisResourceOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r PatchSchedulesListByRedisResourceOperationResponse) LoadMore(ctx context.Context) (resp PatchSchedulesListByRedisResourceOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// PatchSchedulesListByRedisResource ...
func (c RedisClient) PatchSchedulesListByRedisResource(ctx context.Context, id RediId) (resp PatchSchedulesListByRedisResourceOperationResponse, err error) {
	req, err := c.preparerForPatchSchedulesListByRedisResource(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesListByRedisResource", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesListByRedisResource", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForPatchSchedulesListByRedisResource(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesListByRedisResource", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForPatchSchedulesListByRedisResource prepares the PatchSchedulesListByRedisResource request.
func (c RedisClient) preparerForPatchSchedulesListByRedisResource(ctx context.Context, id RediId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/patchSchedules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForPatchSchedulesListByRedisResourceWithNextLink prepares the PatchSchedulesListByRedisResource request with the given nextLink token.
func (c RedisClient) preparerForPatchSchedulesListByRedisResourceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForPatchSchedulesListByRedisResource handles the response to the PatchSchedulesListByRedisResource request. The method always
// closes the http.Response Body.
func (c RedisClient) responderForPatchSchedulesListByRedisResource(resp *http.Response) (result PatchSchedulesListByRedisResourceOperationResponse, err error) {
	type page struct {
		Values   []RedisPatchSchedule `json:"value"`
		NextLink *string              `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result PatchSchedulesListByRedisResourceOperationResponse, err error) {
			req, err := c.preparerForPatchSchedulesListByRedisResourceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesListByRedisResource", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesListByRedisResource", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForPatchSchedulesListByRedisResource(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "redis.RedisClient", "PatchSchedulesListByRedisResource", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// PatchSchedulesListByRedisResourceComplete retrieves all of the results into a single object
func (c RedisClient) PatchSchedulesListByRedisResourceComplete(ctx context.Context, id RediId) (PatchSchedulesListByRedisResourceCompleteResult, error) {
	return c.PatchSchedulesListByRedisResourceCompleteMatchingPredicate(ctx, id, RedisPatchScheduleOperationPredicate{})
}

// PatchSchedulesListByRedisResourceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c RedisClient) PatchSchedulesListByRedisResourceCompleteMatchingPredicate(ctx context.Context, id RediId, predicate RedisPatchScheduleOperationPredicate) (resp PatchSchedulesListByRedisResourceCompleteResult, err error) {
	items := make([]RedisPatchSchedule, 0)

	page, err := c.PatchSchedulesListByRedisResource(ctx, id)
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

	out := PatchSchedulesListByRedisResourceCompleteResult{
		Items: items,
	}
	return out, nil
}
