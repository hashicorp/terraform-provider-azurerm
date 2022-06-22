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

type RemediationsListForResourceOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Remediation

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (RemediationsListForResourceOperationResponse, error)
}

type RemediationsListForResourceCompleteResult struct {
	Items []Remediation
}

func (r RemediationsListForResourceOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r RemediationsListForResourceOperationResponse) LoadMore(ctx context.Context) (resp RemediationsListForResourceOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type RemediationsListForResourceOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultRemediationsListForResourceOperationOptions() RemediationsListForResourceOperationOptions {
	return RemediationsListForResourceOperationOptions{}
}

func (o RemediationsListForResourceOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o RemediationsListForResourceOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// RemediationsListForResource ...
func (c PolicyInsightsClient) RemediationsListForResource(ctx context.Context, id commonids.ScopeId, options RemediationsListForResourceOperationOptions) (resp RemediationsListForResourceOperationResponse, err error) {
	req, err := c.preparerForRemediationsListForResource(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForResource", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForResource", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForRemediationsListForResource(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForResource", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// RemediationsListForResourceComplete retrieves all of the results into a single object
func (c PolicyInsightsClient) RemediationsListForResourceComplete(ctx context.Context, id commonids.ScopeId, options RemediationsListForResourceOperationOptions) (RemediationsListForResourceCompleteResult, error) {
	return c.RemediationsListForResourceCompleteMatchingPredicate(ctx, id, options, RemediationOperationPredicate{})
}

// RemediationsListForResourceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c PolicyInsightsClient) RemediationsListForResourceCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, options RemediationsListForResourceOperationOptions, predicate RemediationOperationPredicate) (resp RemediationsListForResourceCompleteResult, err error) {
	items := make([]Remediation, 0)

	page, err := c.RemediationsListForResource(ctx, id, options)
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

	out := RemediationsListForResourceCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForRemediationsListForResource prepares the RemediationsListForResource request.
func (c PolicyInsightsClient) preparerForRemediationsListForResource(ctx context.Context, id commonids.ScopeId, options RemediationsListForResourceOperationOptions) (*http.Request, error) {
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

// preparerForRemediationsListForResourceWithNextLink prepares the RemediationsListForResource request with the given nextLink token.
func (c PolicyInsightsClient) preparerForRemediationsListForResourceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForRemediationsListForResource handles the response to the RemediationsListForResource request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsListForResource(resp *http.Response) (result RemediationsListForResourceOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result RemediationsListForResourceOperationResponse, err error) {
			req, err := c.preparerForRemediationsListForResourceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForResource", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForResource", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForRemediationsListForResource(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForResource", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
