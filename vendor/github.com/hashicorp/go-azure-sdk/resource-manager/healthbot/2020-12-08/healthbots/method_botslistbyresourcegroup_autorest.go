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

type BotsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]HealthBot

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (BotsListByResourceGroupOperationResponse, error)
}

type BotsListByResourceGroupCompleteResult struct {
	Items []HealthBot
}

func (r BotsListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r BotsListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp BotsListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// BotsListByResourceGroup ...
func (c HealthbotsClient) BotsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp BotsListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForBotsListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForBotsListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForBotsListByResourceGroup prepares the BotsListByResourceGroup request.
func (c HealthbotsClient) preparerForBotsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
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

// preparerForBotsListByResourceGroupWithNextLink prepares the BotsListByResourceGroup request with the given nextLink token.
func (c HealthbotsClient) preparerForBotsListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForBotsListByResourceGroup handles the response to the BotsListByResourceGroup request. The method always
// closes the http.Response Body.
func (c HealthbotsClient) responderForBotsListByResourceGroup(resp *http.Response) (result BotsListByResourceGroupOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result BotsListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForBotsListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForBotsListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "healthbots.HealthbotsClient", "BotsListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// BotsListByResourceGroupComplete retrieves all of the results into a single object
func (c HealthbotsClient) BotsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (BotsListByResourceGroupCompleteResult, error) {
	return c.BotsListByResourceGroupCompleteMatchingPredicate(ctx, id, HealthBotOperationPredicate{})
}

// BotsListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c HealthbotsClient) BotsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate HealthBotOperationPredicate) (resp BotsListByResourceGroupCompleteResult, err error) {
	items := make([]HealthBot, 0)

	page, err := c.BotsListByResourceGroup(ctx, id)
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

	out := BotsListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
