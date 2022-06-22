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

type RemediationsListDeploymentsAtResourceOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]RemediationDeployment

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (RemediationsListDeploymentsAtResourceOperationResponse, error)
}

type RemediationsListDeploymentsAtResourceCompleteResult struct {
	Items []RemediationDeployment
}

func (r RemediationsListDeploymentsAtResourceOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r RemediationsListDeploymentsAtResourceOperationResponse) LoadMore(ctx context.Context) (resp RemediationsListDeploymentsAtResourceOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type RemediationsListDeploymentsAtResourceOperationOptions struct {
	Top *int64
}

func DefaultRemediationsListDeploymentsAtResourceOperationOptions() RemediationsListDeploymentsAtResourceOperationOptions {
	return RemediationsListDeploymentsAtResourceOperationOptions{}
}

func (o RemediationsListDeploymentsAtResourceOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o RemediationsListDeploymentsAtResourceOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// RemediationsListDeploymentsAtResource ...
func (c PolicyInsightsClient) RemediationsListDeploymentsAtResource(ctx context.Context, id ScopedRemediationId, options RemediationsListDeploymentsAtResourceOperationOptions) (resp RemediationsListDeploymentsAtResourceOperationResponse, err error) {
	req, err := c.preparerForRemediationsListDeploymentsAtResource(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtResource", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtResource", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForRemediationsListDeploymentsAtResource(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtResource", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// RemediationsListDeploymentsAtResourceComplete retrieves all of the results into a single object
func (c PolicyInsightsClient) RemediationsListDeploymentsAtResourceComplete(ctx context.Context, id ScopedRemediationId, options RemediationsListDeploymentsAtResourceOperationOptions) (RemediationsListDeploymentsAtResourceCompleteResult, error) {
	return c.RemediationsListDeploymentsAtResourceCompleteMatchingPredicate(ctx, id, options, RemediationDeploymentOperationPredicate{})
}

// RemediationsListDeploymentsAtResourceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c PolicyInsightsClient) RemediationsListDeploymentsAtResourceCompleteMatchingPredicate(ctx context.Context, id ScopedRemediationId, options RemediationsListDeploymentsAtResourceOperationOptions, predicate RemediationDeploymentOperationPredicate) (resp RemediationsListDeploymentsAtResourceCompleteResult, err error) {
	items := make([]RemediationDeployment, 0)

	page, err := c.RemediationsListDeploymentsAtResource(ctx, id, options)
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

	out := RemediationsListDeploymentsAtResourceCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForRemediationsListDeploymentsAtResource prepares the RemediationsListDeploymentsAtResource request.
func (c PolicyInsightsClient) preparerForRemediationsListDeploymentsAtResource(ctx context.Context, id ScopedRemediationId, options RemediationsListDeploymentsAtResourceOperationOptions) (*http.Request, error) {
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

// preparerForRemediationsListDeploymentsAtResourceWithNextLink prepares the RemediationsListDeploymentsAtResource request with the given nextLink token.
func (c PolicyInsightsClient) preparerForRemediationsListDeploymentsAtResourceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForRemediationsListDeploymentsAtResource handles the response to the RemediationsListDeploymentsAtResource request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsListDeploymentsAtResource(resp *http.Response) (result RemediationsListDeploymentsAtResourceOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result RemediationsListDeploymentsAtResourceOperationResponse, err error) {
			req, err := c.preparerForRemediationsListDeploymentsAtResourceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtResource", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtResource", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForRemediationsListDeploymentsAtResource(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtResource", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
