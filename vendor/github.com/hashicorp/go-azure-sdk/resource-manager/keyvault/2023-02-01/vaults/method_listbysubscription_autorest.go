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

type ListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Vault

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListBySubscriptionOperationResponse, error)
}

type ListBySubscriptionCompleteResult struct {
	Items []Vault
}

func (r ListBySubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListBySubscriptionOperationResponse) LoadMore(ctx context.Context) (resp ListBySubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListBySubscriptionOperationOptions struct {
	Top *int64
}

func DefaultListBySubscriptionOperationOptions() ListBySubscriptionOperationOptions {
	return ListBySubscriptionOperationOptions{}
}

func (o ListBySubscriptionOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListBySubscriptionOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListBySubscription ...
func (c VaultsClient) ListBySubscription(ctx context.Context, id commonids.SubscriptionId, options ListBySubscriptionOperationOptions) (resp ListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForListBySubscription(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListBySubscription prepares the ListBySubscription request.
func (c VaultsClient) preparerForListBySubscription(ctx context.Context, id commonids.SubscriptionId, options ListBySubscriptionOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.KeyVault/vaults", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListBySubscriptionWithNextLink prepares the ListBySubscription request with the given nextLink token.
func (c VaultsClient) preparerForListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListBySubscription handles the response to the ListBySubscription request. The method always
// closes the http.Response Body.
func (c VaultsClient) responderForListBySubscription(resp *http.Response) (result ListBySubscriptionOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListBySubscriptionOperationResponse, err error) {
			req, err := c.preparerForListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "vaults.VaultsClient", "ListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListBySubscriptionComplete retrieves all of the results into a single object
func (c VaultsClient) ListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId, options ListBySubscriptionOperationOptions) (ListBySubscriptionCompleteResult, error) {
	return c.ListBySubscriptionCompleteMatchingPredicate(ctx, id, options, VaultOperationPredicate{})
}

// ListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c VaultsClient) ListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ListBySubscriptionOperationOptions, predicate VaultOperationPredicate) (resp ListBySubscriptionCompleteResult, err error) {
	items := make([]Vault, 0)

	page, err := c.ListBySubscription(ctx, id, options)
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

	out := ListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}
