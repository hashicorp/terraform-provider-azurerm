package policyinsights

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

type RemediationsListDeploymentsAtSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]RemediationDeployment

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (RemediationsListDeploymentsAtSubscriptionOperationResponse, error)
}

type RemediationsListDeploymentsAtSubscriptionCompleteResult struct {
	Items []RemediationDeployment
}

func (r RemediationsListDeploymentsAtSubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r RemediationsListDeploymentsAtSubscriptionOperationResponse) LoadMore(ctx context.Context) (resp RemediationsListDeploymentsAtSubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type RemediationsListDeploymentsAtSubscriptionOperationOptions struct {
	Top *int64
}

func DefaultRemediationsListDeploymentsAtSubscriptionOperationOptions() RemediationsListDeploymentsAtSubscriptionOperationOptions {
	return RemediationsListDeploymentsAtSubscriptionOperationOptions{}
}

func (o RemediationsListDeploymentsAtSubscriptionOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o RemediationsListDeploymentsAtSubscriptionOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// RemediationsListDeploymentsAtSubscription ...
func (c PolicyInsightsClient) RemediationsListDeploymentsAtSubscription(ctx context.Context, id RemediationId, options RemediationsListDeploymentsAtSubscriptionOperationOptions) (resp RemediationsListDeploymentsAtSubscriptionOperationResponse, err error) {
	req, err := c.preparerForRemediationsListDeploymentsAtSubscription(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtSubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtSubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForRemediationsListDeploymentsAtSubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtSubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// RemediationsListDeploymentsAtSubscriptionComplete retrieves all of the results into a single object
func (c PolicyInsightsClient) RemediationsListDeploymentsAtSubscriptionComplete(ctx context.Context, id RemediationId, options RemediationsListDeploymentsAtSubscriptionOperationOptions) (RemediationsListDeploymentsAtSubscriptionCompleteResult, error) {
	return c.RemediationsListDeploymentsAtSubscriptionCompleteMatchingPredicate(ctx, id, options, RemediationDeploymentOperationPredicate{})
}

// RemediationsListDeploymentsAtSubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c PolicyInsightsClient) RemediationsListDeploymentsAtSubscriptionCompleteMatchingPredicate(ctx context.Context, id RemediationId, options RemediationsListDeploymentsAtSubscriptionOperationOptions, predicate RemediationDeploymentOperationPredicate) (resp RemediationsListDeploymentsAtSubscriptionCompleteResult, err error) {
	items := make([]RemediationDeployment, 0)

	page, err := c.RemediationsListDeploymentsAtSubscription(ctx, id, options)
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

	out := RemediationsListDeploymentsAtSubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForRemediationsListDeploymentsAtSubscription prepares the RemediationsListDeploymentsAtSubscription request.
func (c PolicyInsightsClient) preparerForRemediationsListDeploymentsAtSubscription(ctx context.Context, id RemediationId, options RemediationsListDeploymentsAtSubscriptionOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/listDeployments", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForRemediationsListDeploymentsAtSubscriptionWithNextLink prepares the RemediationsListDeploymentsAtSubscription request with the given nextLink token.
func (c PolicyInsightsClient) preparerForRemediationsListDeploymentsAtSubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForRemediationsListDeploymentsAtSubscription handles the response to the RemediationsListDeploymentsAtSubscription request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsListDeploymentsAtSubscription(resp *http.Response) (result RemediationsListDeploymentsAtSubscriptionOperationResponse, err error) {
	type page struct {
		Values   []RemediationDeployment `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result RemediationsListDeploymentsAtSubscriptionOperationResponse, err error) {
			req, err := c.preparerForRemediationsListDeploymentsAtSubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtSubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtSubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForRemediationsListDeploymentsAtSubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtSubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
