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

type ActivityLogAlertsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ActivityLogAlertResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ActivityLogAlertsListByResourceGroupOperationResponse, error)
}

type ActivityLogAlertsListByResourceGroupCompleteResult struct {
	Items []ActivityLogAlertResource
}

func (r ActivityLogAlertsListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ActivityLogAlertsListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp ActivityLogAlertsListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ActivityLogAlertsListByResourceGroup ...
func (c ActivityLogAlertsAPIsClient) ActivityLogAlertsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp ActivityLogAlertsListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForActivityLogAlertsListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForActivityLogAlertsListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForActivityLogAlertsListByResourceGroup prepares the ActivityLogAlertsListByResourceGroup request.
func (c ActivityLogAlertsAPIsClient) preparerForActivityLogAlertsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
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

// preparerForActivityLogAlertsListByResourceGroupWithNextLink prepares the ActivityLogAlertsListByResourceGroup request with the given nextLink token.
func (c ActivityLogAlertsAPIsClient) preparerForActivityLogAlertsListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForActivityLogAlertsListByResourceGroup handles the response to the ActivityLogAlertsListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ActivityLogAlertsAPIsClient) responderForActivityLogAlertsListByResourceGroup(resp *http.Response) (result ActivityLogAlertsListByResourceGroupOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ActivityLogAlertsListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForActivityLogAlertsListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForActivityLogAlertsListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "activitylogalertsapis.ActivityLogAlertsAPIsClient", "ActivityLogAlertsListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ActivityLogAlertsListByResourceGroupComplete retrieves all of the results into a single object
func (c ActivityLogAlertsAPIsClient) ActivityLogAlertsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (ActivityLogAlertsListByResourceGroupCompleteResult, error) {
	return c.ActivityLogAlertsListByResourceGroupCompleteMatchingPredicate(ctx, id, ActivityLogAlertResourceOperationPredicate{})
}

// ActivityLogAlertsListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ActivityLogAlertsAPIsClient) ActivityLogAlertsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ActivityLogAlertResourceOperationPredicate) (resp ActivityLogAlertsListByResourceGroupCompleteResult, err error) {
	items := make([]ActivityLogAlertResource, 0)

	page, err := c.ActivityLogAlertsListByResourceGroup(ctx, id)
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

	out := ActivityLogAlertsListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
