package azuretrafficcollectors

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

type BySubscriptionListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AzureTrafficCollector

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (BySubscriptionListOperationResponse, error)
}

type BySubscriptionListCompleteResult struct {
	Items []AzureTrafficCollector
}

func (r BySubscriptionListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r BySubscriptionListOperationResponse) LoadMore(ctx context.Context) (resp BySubscriptionListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// BySubscriptionList ...
func (c AzureTrafficCollectorsClient) BySubscriptionList(ctx context.Context, id commonids.SubscriptionId) (resp BySubscriptionListOperationResponse, err error) {
	req, err := c.preparerForBySubscriptionList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "azuretrafficcollectors.AzureTrafficCollectorsClient", "BySubscriptionList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "azuretrafficcollectors.AzureTrafficCollectorsClient", "BySubscriptionList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForBySubscriptionList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "azuretrafficcollectors.AzureTrafficCollectorsClient", "BySubscriptionList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForBySubscriptionList prepares the BySubscriptionList request.
func (c AzureTrafficCollectorsClient) preparerForBySubscriptionList(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.NetworkFunction/azureTrafficCollectors", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForBySubscriptionListWithNextLink prepares the BySubscriptionList request with the given nextLink token.
func (c AzureTrafficCollectorsClient) preparerForBySubscriptionListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForBySubscriptionList handles the response to the BySubscriptionList request. The method always
// closes the http.Response Body.
func (c AzureTrafficCollectorsClient) responderForBySubscriptionList(resp *http.Response) (result BySubscriptionListOperationResponse, err error) {
	type page struct {
		Values   []AzureTrafficCollector `json:"value"`
		NextLink *string                 `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result BySubscriptionListOperationResponse, err error) {
			req, err := c.preparerForBySubscriptionListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "azuretrafficcollectors.AzureTrafficCollectorsClient", "BySubscriptionList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "azuretrafficcollectors.AzureTrafficCollectorsClient", "BySubscriptionList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForBySubscriptionList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "azuretrafficcollectors.AzureTrafficCollectorsClient", "BySubscriptionList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// BySubscriptionListComplete retrieves all of the results into a single object
func (c AzureTrafficCollectorsClient) BySubscriptionListComplete(ctx context.Context, id commonids.SubscriptionId) (BySubscriptionListCompleteResult, error) {
	return c.BySubscriptionListCompleteMatchingPredicate(ctx, id, AzureTrafficCollectorOperationPredicate{})
}

// BySubscriptionListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AzureTrafficCollectorsClient) BySubscriptionListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate AzureTrafficCollectorOperationPredicate) (resp BySubscriptionListCompleteResult, err error) {
	items := make([]AzureTrafficCollector, 0)

	page, err := c.BySubscriptionList(ctx, id)
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

	out := BySubscriptionListCompleteResult{
		Items: items,
	}
	return out, nil
}
