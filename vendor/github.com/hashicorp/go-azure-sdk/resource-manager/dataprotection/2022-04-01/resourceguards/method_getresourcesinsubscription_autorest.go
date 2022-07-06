package resourceguards

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

type GetResourcesInSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ResourceGuardResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (GetResourcesInSubscriptionOperationResponse, error)
}

type GetResourcesInSubscriptionCompleteResult struct {
	Items []ResourceGuardResource
}

func (r GetResourcesInSubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r GetResourcesInSubscriptionOperationResponse) LoadMore(ctx context.Context) (resp GetResourcesInSubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// GetResourcesInSubscription ...
func (c ResourceGuardsClient) GetResourcesInSubscription(ctx context.Context, id commonids.SubscriptionId) (resp GetResourcesInSubscriptionOperationResponse, err error) {
	req, err := c.preparerForGetResourcesInSubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetResourcesInSubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetResourcesInSubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForGetResourcesInSubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetResourcesInSubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// GetResourcesInSubscriptionComplete retrieves all of the results into a single object
func (c ResourceGuardsClient) GetResourcesInSubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (GetResourcesInSubscriptionCompleteResult, error) {
	return c.GetResourcesInSubscriptionCompleteMatchingPredicate(ctx, id, ResourceGuardResourceOperationPredicate{})
}

// GetResourcesInSubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ResourceGuardsClient) GetResourcesInSubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ResourceGuardResourceOperationPredicate) (resp GetResourcesInSubscriptionCompleteResult, err error) {
	items := make([]ResourceGuardResource, 0)

	page, err := c.GetResourcesInSubscription(ctx, id)
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

	out := GetResourcesInSubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForGetResourcesInSubscription prepares the GetResourcesInSubscription request.
func (c ResourceGuardsClient) preparerForGetResourcesInSubscription(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.DataProtection/resourceGuards", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForGetResourcesInSubscriptionWithNextLink prepares the GetResourcesInSubscription request with the given nextLink token.
func (c ResourceGuardsClient) preparerForGetResourcesInSubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForGetResourcesInSubscription handles the response to the GetResourcesInSubscription request. The method always
// closes the http.Response Body.
func (c ResourceGuardsClient) responderForGetResourcesInSubscription(resp *http.Response) (result GetResourcesInSubscriptionOperationResponse, err error) {
	type page struct {
		Values   []ResourceGuardResource `json:"value"`
		NextLink *string                 `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result GetResourcesInSubscriptionOperationResponse, err error) {
			req, err := c.preparerForGetResourcesInSubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetResourcesInSubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetResourcesInSubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForGetResourcesInSubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetResourcesInSubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
