package policyassignments

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

type ListForResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]PolicyAssignment

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListForResourceGroupOperationResponse, error)
}

type ListForResourceGroupCompleteResult struct {
	Items []PolicyAssignment
}

func (r ListForResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListForResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp ListForResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListForResourceGroupOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListForResourceGroupOperationOptions() ListForResourceGroupOperationOptions {
	return ListForResourceGroupOperationOptions{}
}

func (o ListForResourceGroupOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListForResourceGroupOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListForResourceGroup ...
func (c PolicyAssignmentsClient) ListForResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options ListForResourceGroupOperationOptions) (resp ListForResourceGroupOperationResponse, err error) {
	req, err := c.preparerForListForResourceGroup(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "ListForResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "ListForResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListForResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "ListForResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListForResourceGroup prepares the ListForResourceGroup request.
func (c PolicyAssignmentsClient) preparerForListForResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options ListForResourceGroupOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Authorization/policyAssignments", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListForResourceGroupWithNextLink prepares the ListForResourceGroup request with the given nextLink token.
func (c PolicyAssignmentsClient) preparerForListForResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListForResourceGroup handles the response to the ListForResourceGroup request. The method always
// closes the http.Response Body.
func (c PolicyAssignmentsClient) responderForListForResourceGroup(resp *http.Response) (result ListForResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []PolicyAssignment `json:"value"`
		NextLink *string            `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListForResourceGroupOperationResponse, err error) {
			req, err := c.preparerForListForResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "ListForResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "ListForResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListForResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "policyassignments.PolicyAssignmentsClient", "ListForResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListForResourceGroupComplete retrieves all of the results into a single object
func (c PolicyAssignmentsClient) ListForResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId, options ListForResourceGroupOperationOptions) (ListForResourceGroupCompleteResult, error) {
	return c.ListForResourceGroupCompleteMatchingPredicate(ctx, id, options, PolicyAssignmentOperationPredicate{})
}

// ListForResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c PolicyAssignmentsClient) ListForResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, options ListForResourceGroupOperationOptions, predicate PolicyAssignmentOperationPredicate) (resp ListForResourceGroupCompleteResult, err error) {
	items := make([]PolicyAssignment, 0)

	page, err := c.ListForResourceGroup(ctx, id, options)
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

	out := ListForResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
