package healthbots

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

type BotsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]HealthBot

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (BotsListOperationResponse, error)
}

type BotsListCompleteResult struct {
	Items []HealthBot
}

func (r BotsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r BotsListOperationResponse) LoadMore(ctx context.Context) (resp BotsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// BotsList ...
func (c HealthbotsClient) BotsList(ctx context.Context, id commonids.SubscriptionId) (resp BotsListOperationResponse, err error) {
	req, err := c.preparerForBotsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForBotsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForBotsList prepares the BotsList request.
func (c HealthbotsClient) preparerForBotsList(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.HealthBot/healthBots", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForBotsListWithNextLink prepares the BotsList request with the given nextLink token.
func (c HealthbotsClient) preparerForBotsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForBotsList handles the response to the BotsList request. The method always
// closes the http.Response Body.
func (c HealthbotsClient) responderForBotsList(resp *http.Response) (result BotsListOperationResponse, err error) {
	type page struct {
		Values   []HealthBot `json:"value"`
		NextLink *string     `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result BotsListOperationResponse, err error) {
			req, err := c.preparerForBotsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForBotsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// BotsListComplete retrieves all of the results into a single object
func (c HealthbotsClient) BotsListComplete(ctx context.Context, id commonids.SubscriptionId) (BotsListCompleteResult, error) {
	return c.BotsListCompleteMatchingPredicate(ctx, id, HealthBotOperationPredicate{})
}

// BotsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c HealthbotsClient) BotsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate HealthBotOperationPredicate) (resp BotsListCompleteResult, err error) {
	items := make([]HealthBot, 0)

	page, err := c.BotsList(ctx, id)
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

	out := BotsListCompleteResult{
		Items: items,
	}
	return out, nil
}
