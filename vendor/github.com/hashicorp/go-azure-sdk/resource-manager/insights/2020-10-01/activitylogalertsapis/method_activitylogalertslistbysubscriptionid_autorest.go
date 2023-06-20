package activitylogalertsapis

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

type ActivityLogAlertsListBySubscriptionIdOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ActivityLogAlertResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ActivityLogAlertsListBySubscriptionIdOperationResponse, error)
}

type ActivityLogAlertsListBySubscriptionIdCompleteResult struct {
	Items []ActivityLogAlertResource
}

func (r ActivityLogAlertsListBySubscriptionIdOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ActivityLogAlertsListBySubscriptionIdOperationResponse) LoadMore(ctx context.Context) (resp ActivityLogAlertsListBySubscriptionIdOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ActivityLogAlertsListBySubscriptionId ...
func (c ActivityLogAlertsAPIsClient) ActivityLogAlertsListBySubscriptionId(ctx context.Context, id commonids.SubscriptionId) (resp ActivityLogAlertsListBySubscriptionIdOperationResponse, err error) {
	req, err := c.preparerForActivityLogAlertsListBySubscriptionId(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsListBySubscriptionId", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsListBySubscriptionId", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForActivityLogAlertsListBySubscriptionId(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsListBySubscriptionId", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForActivityLogAlertsListBySubscriptionId prepares the ActivityLogAlertsListBySubscriptionId request.
func (c ActivityLogAlertsAPIsClient) preparerForActivityLogAlertsListBySubscriptionId(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Insights/activityLogAlerts", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForActivityLogAlertsListBySubscriptionIdWithNextLink prepares the ActivityLogAlertsListBySubscriptionId request with the given nextLink token.
func (c ActivityLogAlertsAPIsClient) preparerForActivityLogAlertsListBySubscriptionIdWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForActivityLogAlertsListBySubscriptionId handles the response to the ActivityLogAlertsListBySubscriptionId request. The method always
// closes the http.Response Body.
func (c ActivityLogAlertsAPIsClient) responderForActivityLogAlertsListBySubscriptionId(resp *http.Response) (result ActivityLogAlertsListBySubscriptionIdOperationResponse, err error) {
	type page struct {
		Values   []ActivityLogAlertResource `json:"value"`
		NextLink *string                    `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ActivityLogAlertsListBySubscriptionIdOperationResponse, err error) {
			req, err := c.preparerForActivityLogAlertsListBySubscriptionIdWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsListBySubscriptionId", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsListBySubscriptionId", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForActivityLogAlertsListBySubscriptionId(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsListBySubscriptionId", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ActivityLogAlertsListBySubscriptionIdComplete retrieves all of the results into a single object
func (c ActivityLogAlertsAPIsClient) ActivityLogAlertsListBySubscriptionIdComplete(ctx context.Context, id commonids.SubscriptionId) (ActivityLogAlertsListBySubscriptionIdCompleteResult, error) {
	return c.ActivityLogAlertsListBySubscriptionIdCompleteMatchingPredicate(ctx, id, ActivityLogAlertResourceOperationPredicate{})
}

// ActivityLogAlertsListBySubscriptionIdCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ActivityLogAlertsAPIsClient) ActivityLogAlertsListBySubscriptionIdCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ActivityLogAlertResourceOperationPredicate) (resp ActivityLogAlertsListBySubscriptionIdCompleteResult, err error) {
	items := make([]ActivityLogAlertResource, 0)

	page, err := c.ActivityLogAlertsListBySubscriptionId(ctx, id)
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

	out := ActivityLogAlertsListBySubscriptionIdCompleteResult{
		Items: items,
	}
	return out, nil
}
