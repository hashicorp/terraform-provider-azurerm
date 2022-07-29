package resource

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

type SpatialAnchorsAccountsListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SpatialAnchorsAccount

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (SpatialAnchorsAccountsListBySubscriptionOperationResponse, error)
}

type SpatialAnchorsAccountsListBySubscriptionCompleteResult struct {
	Items []SpatialAnchorsAccount
}

func (r SpatialAnchorsAccountsListBySubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r SpatialAnchorsAccountsListBySubscriptionOperationResponse) LoadMore(ctx context.Context) (resp SpatialAnchorsAccountsListBySubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// SpatialAnchorsAccountsListBySubscription ...
func (c ResourceClient) SpatialAnchorsAccountsListBySubscription(ctx context.Context, id commonids.SubscriptionId) (resp SpatialAnchorsAccountsListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForSpatialAnchorsAccountsListBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForSpatialAnchorsAccountsListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// SpatialAnchorsAccountsListBySubscriptionComplete retrieves all of the results into a single object
func (c ResourceClient) SpatialAnchorsAccountsListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (SpatialAnchorsAccountsListBySubscriptionCompleteResult, error) {
	return c.SpatialAnchorsAccountsListBySubscriptionCompleteMatchingPredicate(ctx, id, SpatialAnchorsAccountOperationPredicate{})
}

// SpatialAnchorsAccountsListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ResourceClient) SpatialAnchorsAccountsListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate SpatialAnchorsAccountOperationPredicate) (resp SpatialAnchorsAccountsListBySubscriptionCompleteResult, err error) {
	items := make([]SpatialAnchorsAccount, 0)

	page, err := c.SpatialAnchorsAccountsListBySubscription(ctx, id)
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

	out := SpatialAnchorsAccountsListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForSpatialAnchorsAccountsListBySubscription prepares the SpatialAnchorsAccountsListBySubscription request.
func (c ResourceClient) preparerForSpatialAnchorsAccountsListBySubscription(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.MixedReality/spatialAnchorsAccounts", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForSpatialAnchorsAccountsListBySubscriptionWithNextLink prepares the SpatialAnchorsAccountsListBySubscription request with the given nextLink token.
func (c ResourceClient) preparerForSpatialAnchorsAccountsListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForSpatialAnchorsAccountsListBySubscription handles the response to the SpatialAnchorsAccountsListBySubscription request. The method always
// closes the http.Response Body.
func (c ResourceClient) responderForSpatialAnchorsAccountsListBySubscription(resp *http.Response) (result SpatialAnchorsAccountsListBySubscriptionOperationResponse, err error) {
	type page struct {
		Values   []SpatialAnchorsAccount `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result SpatialAnchorsAccountsListBySubscriptionOperationResponse, err error) {
			req, err := c.preparerForSpatialAnchorsAccountsListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForSpatialAnchorsAccountsListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resource.ResourceClient", "SpatialAnchorsAccountsListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
