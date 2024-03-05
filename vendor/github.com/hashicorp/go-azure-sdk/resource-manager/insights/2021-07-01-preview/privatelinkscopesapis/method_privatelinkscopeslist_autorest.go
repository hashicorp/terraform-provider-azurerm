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

type PrivateLinkScopesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AzureMonitorPrivateLinkScope

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (PrivateLinkScopesListOperationResponse, error)
}

type PrivateLinkScopesListCompleteResult struct {
	Items []AzureMonitorPrivateLinkScope
}

func (r PrivateLinkScopesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r PrivateLinkScopesListOperationResponse) LoadMore(ctx context.Context) (resp PrivateLinkScopesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// PrivateLinkScopesList ...
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesList(ctx context.Context, id commonids.SubscriptionId) (resp PrivateLinkScopesListOperationResponse, err error) {
	req, err := c.preparerForPrivateLinkScopesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForPrivateLinkScopesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForPrivateLinkScopesList prepares the PrivateLinkScopesList request.
func (c PrivateLinkScopesAPIsClient) preparerForPrivateLinkScopesList(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
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

// preparerForPrivateLinkScopesListWithNextLink prepares the PrivateLinkScopesList request with the given nextLink token.
func (c PrivateLinkScopesAPIsClient) preparerForPrivateLinkScopesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForPrivateLinkScopesList handles the response to the PrivateLinkScopesList request. The method always
// closes the http.Response Body.
func (c PrivateLinkScopesAPIsClient) responderForPrivateLinkScopesList(resp *http.Response) (result PrivateLinkScopesListOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result PrivateLinkScopesListOperationResponse, err error) {
			req, err := c.preparerForPrivateLinkScopesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForPrivateLinkScopesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "privatelinkscopesapis.PrivateLinkScopesAPIsClient", "PrivateLinkScopesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// PrivateLinkScopesListComplete retrieves all of the results into a single object
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesListComplete(ctx context.Context, id commonids.SubscriptionId) (PrivateLinkScopesListCompleteResult, error) {
	return c.PrivateLinkScopesListCompleteMatchingPredicate(ctx, id, AzureMonitorPrivateLinkScopeOperationPredicate{})
}

// PrivateLinkScopesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c PrivateLinkScopesAPIsClient) PrivateLinkScopesListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate AzureMonitorPrivateLinkScopeOperationPredicate) (resp PrivateLinkScopesListCompleteResult, err error) {
	items := make([]AzureMonitorPrivateLinkScope, 0)

	page, err := c.PrivateLinkScopesList(ctx, id)
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

	out := PrivateLinkScopesListCompleteResult{
		Items: items,
	}
	return out, nil
}
