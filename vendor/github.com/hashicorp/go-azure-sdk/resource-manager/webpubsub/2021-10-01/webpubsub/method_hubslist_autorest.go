package webpubsub

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

type HubsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]WebPubSubHub

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (HubsListOperationResponse, error)
}

type HubsListCompleteResult struct {
	Items []WebPubSubHub
}

func (r HubsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r HubsListOperationResponse) LoadMore(ctx context.Context) (resp HubsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// HubsList ...
func (c WebPubSubClient) HubsList(ctx context.Context, id WebPubSubId) (resp HubsListOperationResponse, err error) {
	req, err := c.preparerForHubsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "HubsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "HubsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForHubsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "HubsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForHubsList prepares the HubsList request.
func (c WebPubSubClient) preparerForHubsList(ctx context.Context, id WebPubSubId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/hubs", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForHubsListWithNextLink prepares the HubsList request with the given nextLink token.
func (c WebPubSubClient) preparerForHubsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForHubsList handles the response to the HubsList request. The method always
// closes the http.Response Body.
func (c WebPubSubClient) responderForHubsList(resp *http.Response) (result HubsListOperationResponse, err error) {
	type page struct {
		Values   []WebPubSubHub `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result HubsListOperationResponse, err error) {
			req, err := c.preparerForHubsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "HubsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "HubsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForHubsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "webpubsub.WebPubSubClient", "HubsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// HubsListComplete retrieves all of the results into a single object
func (c WebPubSubClient) HubsListComplete(ctx context.Context, id WebPubSubId) (HubsListCompleteResult, error) {
	return c.HubsListCompleteMatchingPredicate(ctx, id, WebPubSubHubOperationPredicate{})
}

// HubsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c WebPubSubClient) HubsListCompleteMatchingPredicate(ctx context.Context, id WebPubSubId, predicate WebPubSubHubOperationPredicate) (resp HubsListCompleteResult, err error) {
	items := make([]WebPubSubHub, 0)

	page, err := c.HubsList(ctx, id)
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

	out := HubsListCompleteResult{
		Items: items,
	}
	return out, nil
}
