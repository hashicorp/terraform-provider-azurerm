package patchschedules

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

type ListByRedisResourceOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]RedisPatchSchedule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByRedisResourceOperationResponse, error)
}

type ListByRedisResourceCompleteResult struct {
	Items []RedisPatchSchedule
}

func (r ListByRedisResourceOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByRedisResourceOperationResponse) LoadMore(ctx context.Context) (resp ListByRedisResourceOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByRedisResource ...
func (c PatchSchedulesClient) ListByRedisResource(ctx context.Context, id RediId) (resp ListByRedisResourceOperationResponse, err error) {
	req, err := c.preparerForListByRedisResource(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "patchschedules.PatchSchedulesClient", "ListByRedisResource", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "patchschedules.PatchSchedulesClient", "ListByRedisResource", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByRedisResource(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "patchschedules.PatchSchedulesClient", "ListByRedisResource", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByRedisResource prepares the ListByRedisResource request.
func (c PatchSchedulesClient) preparerForListByRedisResource(ctx context.Context, id RediId) (*http.Request, error) {
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

// preparerForListByRedisResourceWithNextLink prepares the ListByRedisResource request with the given nextLink token.
func (c PatchSchedulesClient) preparerForListByRedisResourceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByRedisResource handles the response to the ListByRedisResource request. The method always
// closes the http.Response Body.
func (c PatchSchedulesClient) responderForListByRedisResource(resp *http.Response) (result ListByRedisResourceOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByRedisResourceOperationResponse, err error) {
			req, err := c.preparerForListByRedisResourceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "patchschedules.PatchSchedulesClient", "ListByRedisResource", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "patchschedules.PatchSchedulesClient", "ListByRedisResource", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByRedisResource(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "patchschedules.PatchSchedulesClient", "ListByRedisResource", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByRedisResourceComplete retrieves all of the results into a single object
func (c PatchSchedulesClient) ListByRedisResourceComplete(ctx context.Context, id RediId) (ListByRedisResourceCompleteResult, error) {
	return c.ListByRedisResourceCompleteMatchingPredicate(ctx, id, RedisPatchScheduleOperationPredicate{})
}

// ListByRedisResourceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c PatchSchedulesClient) ListByRedisResourceCompleteMatchingPredicate(ctx context.Context, id RediId, predicate RedisPatchScheduleOperationPredicate) (resp ListByRedisResourceCompleteResult, err error) {
	items := make([]RedisPatchSchedule, 0)

	page, err := c.ListByRedisResource(ctx, id)
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

	out := ListByRedisResourceCompleteResult{
		Items: items,
	}
	return out, nil
}
