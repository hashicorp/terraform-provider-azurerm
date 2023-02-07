package agreements

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

type MarketplaceAgreementsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DatadogAgreementResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (MarketplaceAgreementsListOperationResponse, error)
}

type MarketplaceAgreementsListCompleteResult struct {
	Items []DatadogAgreementResource
}

func (r MarketplaceAgreementsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r MarketplaceAgreementsListOperationResponse) LoadMore(ctx context.Context) (resp MarketplaceAgreementsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// MarketplaceAgreementsList ...
func (c AgreementsClient) MarketplaceAgreementsList(ctx context.Context, id commonids.SubscriptionId) (resp MarketplaceAgreementsListOperationResponse, err error) {
	req, err := c.preparerForMarketplaceAgreementsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "agreements.AgreementsClient", "MarketplaceAgreementsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "agreements.AgreementsClient", "MarketplaceAgreementsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForMarketplaceAgreementsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "agreements.AgreementsClient", "MarketplaceAgreementsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForMarketplaceAgreementsList prepares the MarketplaceAgreementsList request.
func (c AgreementsClient) preparerForMarketplaceAgreementsList(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Datadog/agreements", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForMarketplaceAgreementsListWithNextLink prepares the MarketplaceAgreementsList request with the given nextLink token.
func (c AgreementsClient) preparerForMarketplaceAgreementsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForMarketplaceAgreementsList handles the response to the MarketplaceAgreementsList request. The method always
// closes the http.Response Body.
func (c AgreementsClient) responderForMarketplaceAgreementsList(resp *http.Response) (result MarketplaceAgreementsListOperationResponse, err error) {
	type page struct {
		Values   []DatadogAgreementResource `json:"value"`
		NextLink *string                    `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result MarketplaceAgreementsListOperationResponse, err error) {
			req, err := c.preparerForMarketplaceAgreementsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "agreements.AgreementsClient", "MarketplaceAgreementsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "agreements.AgreementsClient", "MarketplaceAgreementsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForMarketplaceAgreementsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "agreements.AgreementsClient", "MarketplaceAgreementsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// MarketplaceAgreementsListComplete retrieves all of the results into a single object
func (c AgreementsClient) MarketplaceAgreementsListComplete(ctx context.Context, id commonids.SubscriptionId) (MarketplaceAgreementsListCompleteResult, error) {
	return c.MarketplaceAgreementsListCompleteMatchingPredicate(ctx, id, DatadogAgreementResourceOperationPredicate{})
}

// MarketplaceAgreementsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AgreementsClient) MarketplaceAgreementsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate DatadogAgreementResourceOperationPredicate) (resp MarketplaceAgreementsListCompleteResult, err error) {
	items := make([]DatadogAgreementResource, 0)

	page, err := c.MarketplaceAgreementsList(ctx, id)
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

	out := MarketplaceAgreementsListCompleteResult{
		Items: items,
	}
	return out, nil
}
