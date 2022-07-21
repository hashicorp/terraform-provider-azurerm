package netappaccounts

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

type AccountsListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]NetAppAccount

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AccountsListBySubscriptionOperationResponse, error)
}

type AccountsListBySubscriptionCompleteResult struct {
	Items []NetAppAccount
}

func (r AccountsListBySubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AccountsListBySubscriptionOperationResponse) LoadMore(ctx context.Context) (resp AccountsListBySubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// AccountsListBySubscription ...
func (c NetAppAccountsClient) AccountsListBySubscription(ctx context.Context, id commonids.SubscriptionId) (resp AccountsListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForAccountsListBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAccountsListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// AccountsListBySubscriptionComplete retrieves all of the results into a single object
func (c NetAppAccountsClient) AccountsListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (AccountsListBySubscriptionCompleteResult, error) {
	return c.AccountsListBySubscriptionCompleteMatchingPredicate(ctx, id, NetAppAccountOperationPredicate{})
}

// AccountsListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c NetAppAccountsClient) AccountsListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate NetAppAccountOperationPredicate) (resp AccountsListBySubscriptionCompleteResult, err error) {
	items := make([]NetAppAccount, 0)

	page, err := c.AccountsListBySubscription(ctx, id)
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

	out := AccountsListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForAccountsListBySubscription prepares the AccountsListBySubscription request.
func (c NetAppAccountsClient) preparerForAccountsListBySubscription(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.NetApp/netAppAccounts", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForAccountsListBySubscriptionWithNextLink prepares the AccountsListBySubscription request with the given nextLink token.
func (c NetAppAccountsClient) preparerForAccountsListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAccountsListBySubscription handles the response to the AccountsListBySubscription request. The method always
// closes the http.Response Body.
func (c NetAppAccountsClient) responderForAccountsListBySubscription(resp *http.Response) (result AccountsListBySubscriptionOperationResponse, err error) {
	type page struct {
		Values   []NetAppAccount `json:"value"`
		NextLink *string         `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AccountsListBySubscriptionOperationResponse, err error) {
			req, err := c.preparerForAccountsListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAccountsListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
