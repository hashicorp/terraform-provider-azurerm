package privatelinkscopesapis

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

type PrivateLinkScopesListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AzureMonitorPrivateLinkScope

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (PrivateLinkScopesListByResourceGroupOperationResponse, error)
}

type PrivateLinkScopesListByResourceGroupCompleteResult struct {
	Items []AzureMonitorPrivateLinkScope
}

func (r PrivateLinkScopesListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r PrivateLinkScopesListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp PrivateLinkScopesListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// PrivateLinkScopesListByResourceGroup ...
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp PrivateLinkScopesListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForPrivateLinkScopesListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForPrivateLinkScopesListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForPrivateLinkScopesListByResourceGroup prepares the PrivateLinkScopesListByResourceGroup request.
func (c PrivateLinkScopesAPIsClient) preparerForPrivateLinkScopesListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Insights/privateLinkScopes", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForPrivateLinkScopesListByResourceGroupWithNextLink prepares the PrivateLinkScopesListByResourceGroup request with the given nextLink token.
func (c PrivateLinkScopesAPIsClient) preparerForPrivateLinkScopesListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForPrivateLinkScopesListByResourceGroup handles the response to the PrivateLinkScopesListByResourceGroup request. The method always
// closes the http.Response Body.
func (c PrivateLinkScopesAPIsClient) responderForPrivateLinkScopesListByResourceGroup(resp *http.Response) (result PrivateLinkScopesListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []AzureMonitorPrivateLinkScope `json:"value"`
		NextLink *string                        `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result PrivateLinkScopesListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForPrivateLinkScopesListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForPrivateLinkScopesListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// PrivateLinkScopesListByResourceGroupComplete retrieves all of the results into a single object
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (PrivateLinkScopesListByResourceGroupCompleteResult, error) {
	return c.PrivateLinkScopesListByResourceGroupCompleteMatchingPredicate(ctx, id, AzureMonitorPrivateLinkScopeOperationPredicate{})
}

// PrivateLinkScopesListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate AzureMonitorPrivateLinkScopeOperationPredicate) (resp PrivateLinkScopesListByResourceGroupCompleteResult, err error) {
	items := make([]AzureMonitorPrivateLinkScope, 0)

	page, err := c.PrivateLinkScopesListByResourceGroup(ctx, id)
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

	out := PrivateLinkScopesListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
