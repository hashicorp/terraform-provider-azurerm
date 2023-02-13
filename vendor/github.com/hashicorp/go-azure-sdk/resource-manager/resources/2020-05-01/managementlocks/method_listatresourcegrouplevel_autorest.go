package managementlocks

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

type ListAtResourceGroupLevelOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ManagementLockObject

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListAtResourceGroupLevelOperationResponse, error)
}

type ListAtResourceGroupLevelCompleteResult struct {
	Items []ManagementLockObject
}

func (r ListAtResourceGroupLevelOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListAtResourceGroupLevelOperationResponse) LoadMore(ctx context.Context) (resp ListAtResourceGroupLevelOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListAtResourceGroupLevelOperationOptions struct {
	Filter *string
}

func DefaultListAtResourceGroupLevelOperationOptions() ListAtResourceGroupLevelOperationOptions {
	return ListAtResourceGroupLevelOperationOptions{}
}

func (o ListAtResourceGroupLevelOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListAtResourceGroupLevelOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// ListAtResourceGroupLevel ...
func (c ManagementLocksClient) ListAtResourceGroupLevel(ctx context.Context, id commonids.ResourceGroupId, options ListAtResourceGroupLevelOperationOptions) (resp ListAtResourceGroupLevelOperationResponse, err error) {
	req, err := c.preparerForListAtResourceGroupLevel(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "ListAtResourceGroupLevel", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "ListAtResourceGroupLevel", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListAtResourceGroupLevel(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "ListAtResourceGroupLevel", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListAtResourceGroupLevel prepares the ListAtResourceGroupLevel request.
func (c ManagementLocksClient) preparerForListAtResourceGroupLevel(ctx context.Context, id commonids.ResourceGroupId, options ListAtResourceGroupLevelOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Authorization/locks", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListAtResourceGroupLevelWithNextLink prepares the ListAtResourceGroupLevel request with the given nextLink token.
func (c ManagementLocksClient) preparerForListAtResourceGroupLevelWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListAtResourceGroupLevel handles the response to the ListAtResourceGroupLevel request. The method always
// closes the http.Response Body.
func (c ManagementLocksClient) responderForListAtResourceGroupLevel(resp *http.Response) (result ListAtResourceGroupLevelOperationResponse, err error) {
	type page struct {
		Values   []ManagementLockObject `json:"value"`
		NextLink *string                `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListAtResourceGroupLevelOperationResponse, err error) {
			req, err := c.preparerForListAtResourceGroupLevelWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "ListAtResourceGroupLevel", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "ListAtResourceGroupLevel", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListAtResourceGroupLevel(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "ListAtResourceGroupLevel", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListAtResourceGroupLevelComplete retrieves all of the results into a single object
func (c ManagementLocksClient) ListAtResourceGroupLevelComplete(ctx context.Context, id commonids.ResourceGroupId, options ListAtResourceGroupLevelOperationOptions) (ListAtResourceGroupLevelCompleteResult, error) {
	return c.ListAtResourceGroupLevelCompleteMatchingPredicate(ctx, id, options, ManagementLockObjectOperationPredicate{})
}

// ListAtResourceGroupLevelCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ManagementLocksClient) ListAtResourceGroupLevelCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, options ListAtResourceGroupLevelOperationOptions, predicate ManagementLockObjectOperationPredicate) (resp ListAtResourceGroupLevelCompleteResult, err error) {
	items := make([]ManagementLockObject, 0)

	page, err := c.ListAtResourceGroupLevel(ctx, id, options)
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

	out := ListAtResourceGroupLevelCompleteResult{
		Items: items,
	}
	return out, nil
}
