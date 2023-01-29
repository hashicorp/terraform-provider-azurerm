package accounts

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

type MediaservicesListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]MediaService

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (MediaservicesListBySubscriptionOperationResponse, error)
}

type MediaservicesListBySubscriptionCompleteResult struct {
	Items []MediaService
}

func (r MediaservicesListBySubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r MediaservicesListBySubscriptionOperationResponse) LoadMore(ctx context.Context) (resp MediaservicesListBySubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// MediaservicesListBySubscription ...
func (c AccountsClient) MediaservicesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (resp MediaservicesListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForMediaservicesListBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForMediaservicesListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForMediaservicesListBySubscription prepares the MediaservicesListBySubscription request.
func (c AccountsClient) preparerForMediaservicesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Media/mediaServices", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForMediaservicesListBySubscriptionWithNextLink prepares the MediaservicesListBySubscription request with the given nextLink token.
func (c AccountsClient) preparerForMediaservicesListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForMediaservicesListBySubscription handles the response to the MediaservicesListBySubscription request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForMediaservicesListBySubscription(resp *http.Response) (result MediaservicesListBySubscriptionOperationResponse, err error) {
	type page struct {
		Values   []MediaService `json:"value"`
		NextLink *string        `json:"@odata.nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result MediaservicesListBySubscriptionOperationResponse, err error) {
			req, err := c.preparerForMediaservicesListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForMediaservicesListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "MediaservicesListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// MediaservicesListBySubscriptionComplete retrieves all of the results into a single object
func (c AccountsClient) MediaservicesListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (MediaservicesListBySubscriptionCompleteResult, error) {
	return c.MediaservicesListBySubscriptionCompleteMatchingPredicate(ctx, id, MediaServiceOperationPredicate{})
}

// MediaservicesListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AccountsClient) MediaservicesListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate MediaServiceOperationPredicate) (resp MediaservicesListBySubscriptionCompleteResult, err error) {
	items := make([]MediaService, 0)

	page, err := c.MediaservicesListBySubscription(ctx, id)
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

	out := MediaservicesListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}
