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

type AccountsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]NetAppAccount

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AccountsListOperationResponse, error)
}

type AccountsListCompleteResult struct {
	Items []NetAppAccount
}

func (r AccountsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AccountsListOperationResponse) LoadMore(ctx context.Context) (resp AccountsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// AccountsList ...
func (c NetAppAccountsClient) AccountsList(ctx context.Context, id commonids.ResourceGroupId) (resp AccountsListOperationResponse, err error) {
	req, err := c.preparerForAccountsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAccountsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// AccountsListComplete retrieves all of the results into a single object
func (c NetAppAccountsClient) AccountsListComplete(ctx context.Context, id commonids.ResourceGroupId) (AccountsListCompleteResult, error) {
	return c.AccountsListCompleteMatchingPredicate(ctx, id, NetAppAccountOperationPredicate{})
}

// AccountsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c NetAppAccountsClient) AccountsListCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate NetAppAccountOperationPredicate) (resp AccountsListCompleteResult, err error) {
	items := make([]NetAppAccount, 0)

	page, err := c.AccountsList(ctx, id)
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

	out := AccountsListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForAccountsList prepares the AccountsList request.
func (c NetAppAccountsClient) preparerForAccountsList(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
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

// preparerForAccountsListWithNextLink prepares the AccountsList request with the given nextLink token.
func (c NetAppAccountsClient) preparerForAccountsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAccountsList handles the response to the AccountsList request. The method always
// closes the http.Response Body.
func (c NetAppAccountsClient) responderForAccountsList(resp *http.Response) (result AccountsListOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AccountsListOperationResponse, err error) {
			req, err := c.preparerForAccountsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAccountsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "netappaccounts.NetAppAccountsClient", "AccountsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
