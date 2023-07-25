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

type ListAtResourceLevelOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ManagementLockObject

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListAtResourceLevelOperationResponse, error)
}

type ListAtResourceLevelCompleteResult struct {
	Items []ManagementLockObject
}

func (r ListAtResourceLevelOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListAtResourceLevelOperationResponse) LoadMore(ctx context.Context) (resp ListAtResourceLevelOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListAtResourceLevelOperationOptions struct {
	Filter *string
}

func DefaultListAtResourceLevelOperationOptions() ListAtResourceLevelOperationOptions {
	return ListAtResourceLevelOperationOptions{}
}

func (o ListAtResourceLevelOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListAtResourceLevelOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// ListAtResourceLevel ...
func (c ManagementLocksClient) ListAtResourceLevel(ctx context.Context, id commonids.ScopeId, options ListAtResourceLevelOperationOptions) (resp ListAtResourceLevelOperationResponse, err error) {
	req, err := c.preparerForListAtResourceLevel(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "ListAtResourceLevel", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "ListAtResourceLevel", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListAtResourceLevel(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "ListAtResourceLevel", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListAtResourceLevel prepares the ListAtResourceLevel request.
func (c ManagementLocksClient) preparerForListAtResourceLevel(ctx context.Context, id commonids.ScopeId, options ListAtResourceLevelOperationOptions) (*http.Request, error) {
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

// preparerForListAtResourceLevelWithNextLink prepares the ListAtResourceLevel request with the given nextLink token.
func (c ManagementLocksClient) preparerForListAtResourceLevelWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListAtResourceLevel handles the response to the ListAtResourceLevel request. The method always
// closes the http.Response Body.
func (c ManagementLocksClient) responderForListAtResourceLevel(resp *http.Response) (result ListAtResourceLevelOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListAtResourceLevelOperationResponse, err error) {
			req, err := c.preparerForListAtResourceLevelWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "ListAtResourceLevel", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "ListAtResourceLevel", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListAtResourceLevel(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managementlocks.ManagementLocksClient", "ListAtResourceLevel", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListAtResourceLevelComplete retrieves all of the results into a single object
func (c ManagementLocksClient) ListAtResourceLevelComplete(ctx context.Context, id commonids.ScopeId, options ListAtResourceLevelOperationOptions) (ListAtResourceLevelCompleteResult, error) {
	return c.ListAtResourceLevelCompleteMatchingPredicate(ctx, id, options, ManagementLockObjectOperationPredicate{})
}

// ListAtResourceLevelCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ManagementLocksClient) ListAtResourceLevelCompleteMatchingPredicate(ctx context.Context, id commonids.ScopeId, options ListAtResourceLevelOperationOptions, predicate ManagementLockObjectOperationPredicate) (resp ListAtResourceLevelCompleteResult, err error) {
	items := make([]ManagementLockObject, 0)

	page, err := c.ListAtResourceLevel(ctx, id, options)
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

	out := ListAtResourceLevelCompleteResult{
		Items: items,
	}
	return out, nil
}
