package schedules

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

type ListApplicableOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Schedule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListApplicableOperationResponse, error)
}

type ListApplicableCompleteResult struct {
	Items []Schedule
}

func (r ListApplicableOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListApplicableOperationResponse) LoadMore(ctx context.Context) (resp ListApplicableOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListApplicable ...
func (c SchedulesClient) ListApplicable(ctx context.Context, id LabScheduleId) (resp ListApplicableOperationResponse, err error) {
	req, err := c.preparerForListApplicable(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "schedules.SchedulesClient", "ListApplicable", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "schedules.SchedulesClient", "ListApplicable", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListApplicable(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "schedules.SchedulesClient", "ListApplicable", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListApplicable prepares the ListApplicable request.
func (c SchedulesClient) preparerForListApplicable(ctx context.Context, id LabScheduleId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listApplicable", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListApplicableWithNextLink prepares the ListApplicable request with the given nextLink token.
func (c SchedulesClient) preparerForListApplicableWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListApplicable handles the response to the ListApplicable request. The method always
// closes the http.Response Body.
func (c SchedulesClient) responderForListApplicable(resp *http.Response) (result ListApplicableOperationResponse, err error) {
	type page struct {
		Values   []Schedule `json:"value"`
		NextLink *string    `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListApplicableOperationResponse, err error) {
			req, err := c.preparerForListApplicableWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "schedules.SchedulesClient", "ListApplicable", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "schedules.SchedulesClient", "ListApplicable", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListApplicable(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "schedules.SchedulesClient", "ListApplicable", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListApplicableComplete retrieves all of the results into a single object
func (c SchedulesClient) ListApplicableComplete(ctx context.Context, id LabScheduleId) (ListApplicableCompleteResult, error) {
	return c.ListApplicableCompleteMatchingPredicate(ctx, id, ScheduleOperationPredicate{})
}

// ListApplicableCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SchedulesClient) ListApplicableCompleteMatchingPredicate(ctx context.Context, id LabScheduleId, predicate ScheduleOperationPredicate) (resp ListApplicableCompleteResult, err error) {
	items := make([]Schedule, 0)

	page, err := c.ListApplicable(ctx, id)
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

	out := ListApplicableCompleteResult{
		Items: items,
	}
	return out, nil
}
