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

type ListBySubscriptionIdOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Vault

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListBySubscriptionIdOperationResponse, error)
}

type ListBySubscriptionIdCompleteResult struct {
	Items []Vault
}

func (r ListBySubscriptionIdOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListBySubscriptionIdOperationResponse) LoadMore(ctx context.Context) (resp ListBySubscriptionIdOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListBySubscriptionId ...
func (c VaultsClient) ListBySubscriptionId(ctx context.Context, id commonids.SubscriptionId) (resp ListBySubscriptionIdOperationResponse, err error) {
	req, err := c.preparerForListBySubscriptionId(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListBySubscriptionId", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListBySubscriptionId", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListBySubscriptionId(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListBySubscriptionId", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListBySubscriptionId prepares the ListBySubscriptionId request.
func (c VaultsClient) preparerForListBySubscriptionId(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.RecoveryServices/vaults", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListBySubscriptionIdWithNextLink prepares the ListBySubscriptionId request with the given nextLink token.
func (c VaultsClient) preparerForListBySubscriptionIdWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListBySubscriptionId handles the response to the ListBySubscriptionId request. The method always
// closes the http.Response Body.
func (c VaultsClient) responderForListBySubscriptionId(resp *http.Response) (result ListBySubscriptionIdOperationResponse, err error) {
	type page struct {
		Values   []Vault `json:"value"`
		NextLink *string `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListBySubscriptionIdOperationResponse, err error) {
			req, err := c.preparerForListBySubscriptionIdWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListBySubscriptionId", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListBySubscriptionId", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListBySubscriptionId(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListBySubscriptionId", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListBySubscriptionIdComplete retrieves all of the results into a single object
func (c VaultsClient) ListBySubscriptionIdComplete(ctx context.Context, id commonids.SubscriptionId) (ListBySubscriptionIdCompleteResult, error) {
	return c.ListBySubscriptionIdCompleteMatchingPredicate(ctx, id, VaultOperationPredicate{})
}

// ListBySubscriptionIdCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c VaultsClient) ListBySubscriptionIdCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate VaultOperationPredicate) (resp ListBySubscriptionIdCompleteResult, err error) {
	items := make([]Vault, 0)

	page, err := c.ListBySubscriptionId(ctx, id)
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

	out := ListBySubscriptionIdCompleteResult{
		Items: items,
	}
	return out, nil
}
