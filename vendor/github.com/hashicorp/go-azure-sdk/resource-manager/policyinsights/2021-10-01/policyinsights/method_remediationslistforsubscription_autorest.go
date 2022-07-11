package policyinsights

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

type RemediationsListForSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Remediation

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (RemediationsListForSubscriptionOperationResponse, error)
}

type RemediationsListForSubscriptionCompleteResult struct {
	Items []Remediation
}

func (r RemediationsListForSubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r RemediationsListForSubscriptionOperationResponse) LoadMore(ctx context.Context) (resp RemediationsListForSubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type RemediationsListForSubscriptionOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultRemediationsListForSubscriptionOperationOptions() RemediationsListForSubscriptionOperationOptions {
	return RemediationsListForSubscriptionOperationOptions{}
}

func (o RemediationsListForSubscriptionOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o RemediationsListForSubscriptionOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// RemediationsListForSubscription ...
func (c PolicyInsightsClient) RemediationsListForSubscription(ctx context.Context, id commonids.SubscriptionId, options RemediationsListForSubscriptionOperationOptions) (resp RemediationsListForSubscriptionOperationResponse, err error) {
	req, err := c.preparerForRemediationsListForSubscription(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForSubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForSubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForRemediationsListForSubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForSubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// RemediationsListForSubscriptionComplete retrieves all of the results into a single object
func (c PolicyInsightsClient) RemediationsListForSubscriptionComplete(ctx context.Context, id commonids.SubscriptionId, options RemediationsListForSubscriptionOperationOptions) (RemediationsListForSubscriptionCompleteResult, error) {
	return c.RemediationsListForSubscriptionCompleteMatchingPredicate(ctx, id, options, RemediationOperationPredicate{})
}

// RemediationsListForSubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c PolicyInsightsClient) RemediationsListForSubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options RemediationsListForSubscriptionOperationOptions, predicate RemediationOperationPredicate) (resp RemediationsListForSubscriptionCompleteResult, err error) {
	items := make([]Remediation, 0)

	page, err := c.RemediationsListForSubscription(ctx, id, options)
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

	out := RemediationsListForSubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForRemediationsListForSubscription prepares the RemediationsListForSubscription request.
func (c PolicyInsightsClient) preparerForRemediationsListForSubscription(ctx context.Context, id commonids.SubscriptionId, options RemediationsListForSubscriptionOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.PolicyInsights/remediations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForRemediationsListForSubscriptionWithNextLink prepares the RemediationsListForSubscription request with the given nextLink token.
func (c PolicyInsightsClient) preparerForRemediationsListForSubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForRemediationsListForSubscription handles the response to the RemediationsListForSubscription request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsListForSubscription(resp *http.Response) (result RemediationsListForSubscriptionOperationResponse, err error) {
	type page struct {
		Values   []Remediation `json:"value"`
		NextLink *string       `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result RemediationsListForSubscriptionOperationResponse, err error) {
			req, err := c.preparerForRemediationsListForSubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForSubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForSubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForRemediationsListForSubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForSubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
