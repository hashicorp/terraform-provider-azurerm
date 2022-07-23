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

type RemediationsListDeploymentsAtResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]RemediationDeployment

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (RemediationsListDeploymentsAtResourceGroupOperationResponse, error)
}

type RemediationsListDeploymentsAtResourceGroupCompleteResult struct {
	Items []RemediationDeployment
}

func (r RemediationsListDeploymentsAtResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r RemediationsListDeploymentsAtResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp RemediationsListDeploymentsAtResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type RemediationsListDeploymentsAtResourceGroupOperationOptions struct {
	Top *int64
}

func DefaultRemediationsListDeploymentsAtResourceGroupOperationOptions() RemediationsListDeploymentsAtResourceGroupOperationOptions {
	return RemediationsListDeploymentsAtResourceGroupOperationOptions{}
}

func (o RemediationsListDeploymentsAtResourceGroupOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o RemediationsListDeploymentsAtResourceGroupOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// RemediationsListDeploymentsAtResourceGroup ...
func (c PolicyInsightsClient) RemediationsListDeploymentsAtResourceGroup(ctx context.Context, id ProviderRemediationId, options RemediationsListDeploymentsAtResourceGroupOperationOptions) (resp RemediationsListDeploymentsAtResourceGroupOperationResponse, err error) {
	req, err := c.preparerForRemediationsListDeploymentsAtResourceGroup(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForRemediationsListDeploymentsAtResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// RemediationsListDeploymentsAtResourceGroupComplete retrieves all of the results into a single object
func (c PolicyInsightsClient) RemediationsListDeploymentsAtResourceGroupComplete(ctx context.Context, id ProviderRemediationId, options RemediationsListDeploymentsAtResourceGroupOperationOptions) (RemediationsListDeploymentsAtResourceGroupCompleteResult, error) {
	return c.RemediationsListDeploymentsAtResourceGroupCompleteMatchingPredicate(ctx, id, options, RemediationDeploymentOperationPredicate{})
}

// RemediationsListDeploymentsAtResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c PolicyInsightsClient) RemediationsListDeploymentsAtResourceGroupCompleteMatchingPredicate(ctx context.Context, id ProviderRemediationId, options RemediationsListDeploymentsAtResourceGroupOperationOptions, predicate RemediationDeploymentOperationPredicate) (resp RemediationsListDeploymentsAtResourceGroupCompleteResult, err error) {
	items := make([]RemediationDeployment, 0)

	page, err := c.RemediationsListDeploymentsAtResourceGroup(ctx, id, options)
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

	out := RemediationsListDeploymentsAtResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForRemediationsListDeploymentsAtResourceGroup prepares the RemediationsListDeploymentsAtResourceGroup request.
func (c PolicyInsightsClient) preparerForRemediationsListDeploymentsAtResourceGroup(ctx context.Context, id ProviderRemediationId, options RemediationsListDeploymentsAtResourceGroupOperationOptions) (*http.Request, error) {
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

// preparerForRemediationsListDeploymentsAtResourceGroupWithNextLink prepares the RemediationsListDeploymentsAtResourceGroup request with the given nextLink token.
func (c PolicyInsightsClient) preparerForRemediationsListDeploymentsAtResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForRemediationsListDeploymentsAtResourceGroup handles the response to the RemediationsListDeploymentsAtResourceGroup request. The method always
// closes the http.Response Body.
func (c PolicyInsightsClient) responderForRemediationsListDeploymentsAtResourceGroup(resp *http.Response) (result RemediationsListDeploymentsAtResourceGroupOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result RemediationsListDeploymentsAtResourceGroupOperationResponse, err error) {
			req, err := c.preparerForRemediationsListDeploymentsAtResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForRemediationsListDeploymentsAtResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyinsights.PolicyInsightsClient", "RemediationsListDeploymentsAtResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
