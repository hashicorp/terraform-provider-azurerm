package account

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
	Model        *[]Account

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListBySubscriptionOperationResponse, error)
}

type ListBySubscriptionCompleteResult struct {
	Items []Account
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

// ListBySubscription ...
func (c AccountClient) ListBySubscription(ctx context.Context, id commonids.SubscriptionId) (resp ListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForListBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.AccountClient", "ListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.AccountClient", "ListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "account.AccountClient", "ListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListBySubscription prepares the ListBySubscription request.
func (c AccountClient) preparerForListBySubscription(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.DataShare/accounts", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListBySubscriptionWithNextLink prepares the ListBySubscription request with the given nextLink token.
func (c AccountClient) preparerForListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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
func (c AccountClient) responderForListBySubscription(resp *http.Response) (result ListBySubscriptionOperationResponse, err error) {
	type page struct {
		Values   []Account `json:"value"`
		NextLink *string   `json:"nextLink"`
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
				err = autorest.NewErrorWithError(err, "account.AccountClient", "ListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "account.AccountClient", "ListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "account.AccountClient", "ListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListBySubscriptionComplete retrieves all of the results into a single object
func (c AccountClient) ListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (ListBySubscriptionCompleteResult, error) {
	return c.ListBySubscriptionCompleteMatchingPredicate(ctx, id, AccountOperationPredicate{})
}

// ListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AccountClient) ListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate AccountOperationPredicate) (resp ListBySubscriptionCompleteResult, err error) {
	items := make([]Account, 0)

	page, err := c.ListBySubscription(ctx, id)
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
