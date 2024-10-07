package vaults

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

type ListDeletedOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DeletedVault

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListDeletedOperationResponse, error)
}

type ListDeletedCompleteResult struct {
	Items []DeletedVault
}

func (r ListDeletedOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListDeletedOperationResponse) LoadMore(ctx context.Context) (resp ListDeletedOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListDeleted ...
func (c VaultsClient) ListDeleted(ctx context.Context, id commonids.SubscriptionId) (resp ListDeletedOperationResponse, err error) {
	req, err := c.preparerForListDeleted(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListDeleted", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListDeleted", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListDeleted(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListDeleted", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListDeleted prepares the ListDeleted request.
func (c VaultsClient) preparerForListDeleted(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.KeyVault/deletedVaults", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListDeletedWithNextLink prepares the ListDeleted request with the given nextLink token.
func (c VaultsClient) preparerForListDeletedWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListDeleted handles the response to the ListDeleted request. The method always
// closes the http.Response Body.
func (c VaultsClient) responderForListDeleted(resp *http.Response) (result ListDeletedOperationResponse, err error) {
	type page struct {
		Values   []DeletedVault `json:"value"`
		NextLink *string        `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListDeletedOperationResponse, err error) {
			req, err := c.preparerForListDeletedWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListDeleted", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListDeleted", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListDeleted(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListDeleted", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListDeletedComplete retrieves all of the results into a single object
func (c VaultsClient) ListDeletedComplete(ctx context.Context, id commonids.SubscriptionId) (ListDeletedCompleteResult, error) {
	return c.ListDeletedCompleteMatchingPredicate(ctx, id, DeletedVaultOperationPredicate{})
}

// ListDeletedCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c VaultsClient) ListDeletedCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate DeletedVaultOperationPredicate) (resp ListDeletedCompleteResult, err error) {
	items := make([]DeletedVault, 0)

	page, err := c.ListDeleted(ctx, id)
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

	out := ListDeletedCompleteResult{
		Items: items,
	}
	return out, nil
}
