package applicationinsights

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbooksListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Workbook

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (WorkbooksListBySubscriptionOperationResponse, error)
}

type WorkbooksListBySubscriptionCompleteResult struct {
	Items []Workbook
}

func (r WorkbooksListBySubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r WorkbooksListBySubscriptionOperationResponse) LoadMore(ctx context.Context) (resp WorkbooksListBySubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type WorkbooksListBySubscriptionOperationOptions struct {
	CanFetchContent *bool
	Category        *CategoryType
	Tags            *string
}

func DefaultWorkbooksListBySubscriptionOperationOptions() WorkbooksListBySubscriptionOperationOptions {
	return WorkbooksListBySubscriptionOperationOptions{}
}

func (o WorkbooksListBySubscriptionOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o WorkbooksListBySubscriptionOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.CanFetchContent != nil {
		out["canFetchContent"] = *o.CanFetchContent
	}

	if o.Category != nil {
		out["category"] = *o.Category
	}

	if o.Tags != nil {
		out["tags"] = *o.Tags
	}

	return out
}

// WorkbooksListBySubscription ...
func (c ApplicationInsightsClient) WorkbooksListBySubscription(ctx context.Context, id commonids.SubscriptionId, options WorkbooksListBySubscriptionOperationOptions) (resp WorkbooksListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForWorkbooksListBySubscription(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForWorkbooksListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// WorkbooksListBySubscriptionComplete retrieves all of the results into a single object
func (c ApplicationInsightsClient) WorkbooksListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId, options WorkbooksListBySubscriptionOperationOptions) (WorkbooksListBySubscriptionCompleteResult, error) {
	return c.WorkbooksListBySubscriptionCompleteMatchingPredicate(ctx, id, options, WorkbookOperationPredicate{})
}

// WorkbooksListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ApplicationInsightsClient) WorkbooksListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options WorkbooksListBySubscriptionOperationOptions, predicate WorkbookOperationPredicate) (resp WorkbooksListBySubscriptionCompleteResult, err error) {
	items := make([]Workbook, 0)

	page, err := c.WorkbooksListBySubscription(ctx, id, options)
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

	out := WorkbooksListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForWorkbooksListBySubscription prepares the WorkbooksListBySubscription request.
func (c ApplicationInsightsClient) preparerForWorkbooksListBySubscription(ctx context.Context, id commonids.SubscriptionId, options WorkbooksListBySubscriptionOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Insights/workbooks", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForWorkbooksListBySubscriptionWithNextLink prepares the WorkbooksListBySubscription request with the given nextLink token.
func (c ApplicationInsightsClient) preparerForWorkbooksListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForWorkbooksListBySubscription handles the response to the WorkbooksListBySubscription request. The method always
// closes the http.Response Body.
func (c ApplicationInsightsClient) responderForWorkbooksListBySubscription(resp *http.Response) (result WorkbooksListBySubscriptionOperationResponse, err error) {
	type page struct {
		Values   []Workbook `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result WorkbooksListBySubscriptionOperationResponse, err error) {
			req, err := c.preparerForWorkbooksListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForWorkbooksListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "applicationinsights.ApplicationInsightsClient", "WorkbooksListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
