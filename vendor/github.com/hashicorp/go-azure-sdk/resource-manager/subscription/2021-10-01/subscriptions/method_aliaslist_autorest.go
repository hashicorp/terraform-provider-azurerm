package subscriptions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AliasListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SubscriptionAliasResponse

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AliasListOperationResponse, error)
}

type AliasListCompleteResult struct {
	Items []SubscriptionAliasResponse
}

func (r AliasListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AliasListOperationResponse) LoadMore(ctx context.Context) (resp AliasListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// AliasList ...
func (c SubscriptionsClient) AliasList(ctx context.Context) (resp AliasListOperationResponse, err error) {
	req, err := c.preparerForAliasList(ctx)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAliasList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForAliasList prepares the AliasList request.
func (c SubscriptionsClient) preparerForAliasList(ctx context.Context) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath("/providers/Microsoft.Subscription/aliases"),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForAliasListWithNextLink prepares the AliasList request with the given nextLink token.
func (c SubscriptionsClient) preparerForAliasListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAliasList handles the response to the AliasList request. The method always
// closes the http.Response Body.
func (c SubscriptionsClient) responderForAliasList(resp *http.Response) (result AliasListOperationResponse, err error) {
	type page struct {
		Values   []SubscriptionAliasResponse `json:"value"`
		NextLink *string                     `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AliasListOperationResponse, err error) {
			req, err := c.preparerForAliasListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAliasList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "subscriptions.SubscriptionsClient", "AliasList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// AliasListComplete retrieves all of the results into a single object
func (c SubscriptionsClient) AliasListComplete(ctx context.Context) (AliasListCompleteResult, error) {
	return c.AliasListCompleteMatchingPredicate(ctx, SubscriptionAliasResponseOperationPredicate{})
}

// AliasListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SubscriptionsClient) AliasListCompleteMatchingPredicate(ctx context.Context, predicate SubscriptionAliasResponseOperationPredicate) (resp AliasListCompleteResult, err error) {
	items := make([]SubscriptionAliasResponse, 0)

	page, err := c.AliasList(ctx)
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

	out := AliasListCompleteResult{
		Items: items,
	}
	return out, nil
}
