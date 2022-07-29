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

type RemediationsListDeploymentsAtManagementGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]RemediationDeployment

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (RemediationsListDeploymentsAtManagementGroupOperationResponse, error)
}

type RemediationsListDeploymentsAtManagementGroupCompleteResult struct {
	Items []RemediationDeployment
}

func (r RemediationsListDeploymentsAtManagementGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r RemediationsListDeploymentsAtManagementGroupOperationResponse) LoadMore(ctx context.Context) (resp RemediationsListDeploymentsAtManagementGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type RemediationsListDeploymentsAtManagementGroupOperationOptions struct {
	Top *int64
}

func DefaultRemediationsListDeploymentsAtManagementGroupOperationOptions() RemediationsListDeploymentsAtManagementGroupOperationOptions {
	return RemediationsListDeploymentsAtManagementGroupOperationOptions{}
}

func (o RemediationsListDeploymentsAtManagementGroupOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o RemediationsListDeploymentsAtManagementGroupOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// RemediationsListDeploymentsAtManagementGroup ...
func (c PolicyInsightsClient) RemediationsListDeploymentsAtManagementGroup(ctx context.Context, id Providers2RemediationId, options RemediationsListDeploymentsAtManagementGroupOperationOptions) (resp RemediationsListDeploymentsAtManagementGroupOperationResponse, err error) {
	req, err := c.preparerForRemediationsListDeploymentsAtManagementGroup(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtManagementGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtManagementGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForRemediationsListDeploymentsAtManagementGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtManagementGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// RemediationsListDeploymentsAtManagementGroupComplete retrieves all of the results into a single object
func (c PolicyInsightsClient) RemediationsListDeploymentsAtManagementGroupComplete(ctx context.Context, id Providers2RemediationId, options RemediationsListDeploymentsAtManagementGroupOperationOptions) (RemediationsListDeploymentsAtManagementGroupCompleteResult, error) {
	return c.RemediationsListDeploymentsAtManagementGroupCompleteMatchingPredicate(ctx, id, options, RemediationDeploymentOperationPredicate{})
}

// RemediationsListDeploymentsAtManagementGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c PolicyInsightsClient) RemediationsListDeploymentsAtManagementGroupCompleteMatchingPredicate(ctx context.Context, id Providers2RemediationId, options RemediationsListDeploymentsAtManagementGroupOperationOptions, predicate RemediationDeploymentOperationPredicate) (resp RemediationsListDeploymentsAtManagementGroupCompleteResult, err error) {
	items := make([]RemediationDeployment, 0)

	page, err := c.RemediationsListDeploymentsAtManagementGroup(ctx, id, options)
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

	out := RemediationsListDeploymentsAtManagementGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForRemediationsListDeploymentsAtManagementGroup prepares the RemediationsListDeploymentsAtManagementGroup request.
func (c PolicyInsightsClient) preparerForRemediationsListDeploymentsAtManagementGroup(ctx context.Context, id Providers2RemediationId, options RemediationsListDeploymentsAtManagementGroupOperationOptions) (*http.Request, error) {
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

// preparerForRemediationsListDeploymentsAtManagementGroupWithNextLink prepares the RemediationsListDeploymentsAtManagementGroup request with the given nextLink token.
func (c PolicyInsightsClient) preparerForRemediationsListDeploymentsAtManagementGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForRemediationsListDeploymentsAtManagementGroup handles the response to the RemediationsListDeploymentsAtManagementGroup request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsListDeploymentsAtManagementGroup(resp *http.Response) (result RemediationsListDeploymentsAtManagementGroupOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result RemediationsListDeploymentsAtManagementGroupOperationResponse, err error) {
			req, err := c.preparerForRemediationsListDeploymentsAtManagementGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtManagementGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtManagementGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForRemediationsListDeploymentsAtManagementGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtManagementGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
