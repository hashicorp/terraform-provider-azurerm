package webapplicationfirewallpolicies

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

type PoliciesListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]WebApplicationFirewallPolicy

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (PoliciesListBySubscriptionOperationResponse, error)
}

type PoliciesListBySubscriptionCompleteResult struct {
	Items []WebApplicationFirewallPolicy
}

func (r PoliciesListBySubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r PoliciesListBySubscriptionOperationResponse) LoadMore(ctx context.Context) (resp PoliciesListBySubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// PoliciesListBySubscription ...
func (c WebApplicationFirewallPoliciesClient) PoliciesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (resp PoliciesListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForPoliciesListBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForPoliciesListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForPoliciesListBySubscription prepares the PoliciesListBySubscription request.
func (c WebApplicationFirewallPoliciesClient) preparerForPoliciesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForPoliciesListBySubscriptionWithNextLink prepares the PoliciesListBySubscription request with the given nextLink token.
func (c WebApplicationFirewallPoliciesClient) preparerForPoliciesListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForPoliciesListBySubscription handles the response to the PoliciesListBySubscription request. The method always
// closes the http.Response Body.
func (c WebApplicationFirewallPoliciesClient) responderForPoliciesListBySubscription(resp *http.Response) (result PoliciesListBySubscriptionOperationResponse, err error) {
	type page struct {
		Values   []WebApplicationFirewallPolicy `json:"value"`
		NextLink *string                        `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result PoliciesListBySubscriptionOperationResponse, err error) {
			req, err := c.preparerForPoliciesListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForPoliciesListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "webapplicationfirewallpolicies.WebApplicationFirewallPoliciesClient", "PoliciesListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// PoliciesListBySubscriptionComplete retrieves all of the results into a single object
func (c WebApplicationFirewallPoliciesClient) PoliciesListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (PoliciesListBySubscriptionCompleteResult, error) {
	return c.PoliciesListBySubscriptionCompleteMatchingPredicate(ctx, id, WebApplicationFirewallPolicyOperationPredicate{})
}

// PoliciesListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c WebApplicationFirewallPoliciesClient) PoliciesListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate WebApplicationFirewallPolicyOperationPredicate) (resp PoliciesListBySubscriptionCompleteResult, err error) {
	items := make([]WebApplicationFirewallPolicy, 0)

	page, err := c.PoliciesListBySubscription(ctx, id)
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

	out := PoliciesListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}
