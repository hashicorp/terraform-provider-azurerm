package share

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

type ProviderShareSubscriptionsListByShareOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ProviderShareSubscription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ProviderShareSubscriptionsListByShareOperationResponse, error)
}

type ProviderShareSubscriptionsListByShareCompleteResult struct {
	Items []ProviderShareSubscription
}

func (r ProviderShareSubscriptionsListByShareOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ProviderShareSubscriptionsListByShareOperationResponse) LoadMore(ctx context.Context) (resp ProviderShareSubscriptionsListByShareOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ProviderShareSubscriptionsListByShare ...
func (c ShareClient) ProviderShareSubscriptionsListByShare(ctx context.Context, id ShareId) (resp ProviderShareSubscriptionsListByShareOperationResponse, err error) {
	req, err := c.preparerForProviderShareSubscriptionsListByShare(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ProviderShareSubscriptionsListByShare", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ProviderShareSubscriptionsListByShare", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForProviderShareSubscriptionsListByShare(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "share.ShareClient", "ProviderShareSubscriptionsListByShare", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForProviderShareSubscriptionsListByShare prepares the ProviderShareSubscriptionsListByShare request.
func (c ShareClient) preparerForProviderShareSubscriptionsListByShare(ctx context.Context, id ShareId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providerShareSubscriptions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForProviderShareSubscriptionsListByShareWithNextLink prepares the ProviderShareSubscriptionsListByShare request with the given nextLink token.
func (c ShareClient) preparerForProviderShareSubscriptionsListByShareWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForProviderShareSubscriptionsListByShare handles the response to the ProviderShareSubscriptionsListByShare request. The method always
// closes the http.Response Body.
func (c ShareClient) responderForProviderShareSubscriptionsListByShare(resp *http.Response) (result ProviderShareSubscriptionsListByShareOperationResponse, err error) {
	type page struct {
		Values   []ProviderShareSubscription `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ProviderShareSubscriptionsListByShareOperationResponse, err error) {
			req, err := c.preparerForProviderShareSubscriptionsListByShareWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "share.ShareClient", "ProviderShareSubscriptionsListByShare", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "share.ShareClient", "ProviderShareSubscriptionsListByShare", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForProviderShareSubscriptionsListByShare(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "share.ShareClient", "ProviderShareSubscriptionsListByShare", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ProviderShareSubscriptionsListByShareComplete retrieves all of the results into a single object
func (c ShareClient) ProviderShareSubscriptionsListByShareComplete(ctx context.Context, id ShareId) (ProviderShareSubscriptionsListByShareCompleteResult, error) {
	return c.ProviderShareSubscriptionsListByShareCompleteMatchingPredicate(ctx, id, ProviderShareSubscriptionOperationPredicate{})
}

// ProviderShareSubscriptionsListByShareCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ShareClient) ProviderShareSubscriptionsListByShareCompleteMatchingPredicate(ctx context.Context, id ShareId, predicate ProviderShareSubscriptionOperationPredicate) (resp ProviderShareSubscriptionsListByShareCompleteResult, err error) {
	items := make([]ProviderShareSubscription, 0)

	page, err := c.ProviderShareSubscriptionsListByShare(ctx, id)
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

	out := ProviderShareSubscriptionsListByShareCompleteResult{
		Items: items,
	}
	return out, nil
}
