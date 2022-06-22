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

type RemediationsListForManagementGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Remediation

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (RemediationsListForManagementGroupOperationResponse, error)
}

type RemediationsListForManagementGroupCompleteResult struct {
	Items []Remediation
}

func (r RemediationsListForManagementGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r RemediationsListForManagementGroupOperationResponse) LoadMore(ctx context.Context) (resp RemediationsListForManagementGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type RemediationsListForManagementGroupOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultRemediationsListForManagementGroupOperationOptions() RemediationsListForManagementGroupOperationOptions {
	return RemediationsListForManagementGroupOperationOptions{}
}

func (o RemediationsListForManagementGroupOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o RemediationsListForManagementGroupOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// RemediationsListForManagementGroup ...
func (c PolicyInsightsClient) RemediationsListForManagementGroup(ctx context.Context, id ManagementGroupId, options RemediationsListForManagementGroupOperationOptions) (resp RemediationsListForManagementGroupOperationResponse, err error) {
	req, err := c.preparerForRemediationsListForManagementGroup(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForManagementGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForManagementGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForRemediationsListForManagementGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForManagementGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// RemediationsListForManagementGroupComplete retrieves all of the results into a single object
func (c PolicyInsightsClient) RemediationsListForManagementGroupComplete(ctx context.Context, id ManagementGroupId, options RemediationsListForManagementGroupOperationOptions) (RemediationsListForManagementGroupCompleteResult, error) {
	return c.RemediationsListForManagementGroupCompleteMatchingPredicate(ctx, id, options, RemediationOperationPredicate{})
}

// RemediationsListForManagementGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c PolicyInsightsClient) RemediationsListForManagementGroupCompleteMatchingPredicate(ctx context.Context, id ManagementGroupId, options RemediationsListForManagementGroupOperationOptions, predicate RemediationOperationPredicate) (resp RemediationsListForManagementGroupCompleteResult, err error) {
	items := make([]Remediation, 0)

	page, err := c.RemediationsListForManagementGroup(ctx, id, options)
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

	out := RemediationsListForManagementGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForRemediationsListForManagementGroup prepares the RemediationsListForManagementGroup request.
func (c PolicyInsightsClient) preparerForRemediationsListForManagementGroup(ctx context.Context, id ManagementGroupId, options RemediationsListForManagementGroupOperationOptions) (*http.Request, error) {
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

// preparerForRemediationsListForManagementGroupWithNextLink prepares the RemediationsListForManagementGroup request with the given nextLink token.
func (c PolicyInsightsClient) preparerForRemediationsListForManagementGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForRemediationsListForManagementGroup handles the response to the RemediationsListForManagementGroup request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsListForManagementGroup(resp *http.Response) (result RemediationsListForManagementGroupOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result RemediationsListForManagementGroupOperationResponse, err error) {
			req, err := c.preparerForRemediationsListForManagementGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForManagementGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForManagementGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForRemediationsListForManagementGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListForManagementGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
