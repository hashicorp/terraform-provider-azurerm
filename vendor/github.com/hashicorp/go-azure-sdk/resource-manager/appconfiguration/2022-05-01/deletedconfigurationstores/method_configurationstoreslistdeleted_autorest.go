package deletedconfigurationstores

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

type ConfigurationStoresListDeletedOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DeletedConfigurationStore

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ConfigurationStoresListDeletedOperationResponse, error)
}

type ConfigurationStoresListDeletedCompleteResult struct {
	Items []DeletedConfigurationStore
}

func (r ConfigurationStoresListDeletedOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ConfigurationStoresListDeletedOperationResponse) LoadMore(ctx context.Context) (resp ConfigurationStoresListDeletedOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ConfigurationStoresListDeleted ...
func (c DeletedConfigurationStoresClient) ConfigurationStoresListDeleted(ctx context.Context, id commonids.SubscriptionId) (resp ConfigurationStoresListDeletedOperationResponse, err error) {
	req, err := c.preparerForConfigurationStoresListDeleted(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deletedconfigurationstores.DeletedConfigurationStoresClient", "ConfigurationStoresListDeleted", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "deletedconfigurationstores.DeletedConfigurationStoresClient", "ConfigurationStoresListDeleted", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForConfigurationStoresListDeleted(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deletedconfigurationstores.DeletedConfigurationStoresClient", "ConfigurationStoresListDeleted", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForConfigurationStoresListDeleted prepares the ConfigurationStoresListDeleted request.
func (c DeletedConfigurationStoresClient) preparerForConfigurationStoresListDeleted(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.AppConfiguration/deletedConfigurationStores", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForConfigurationStoresListDeletedWithNextLink prepares the ConfigurationStoresListDeleted request with the given nextLink token.
func (c DeletedConfigurationStoresClient) preparerForConfigurationStoresListDeletedWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForConfigurationStoresListDeleted handles the response to the ConfigurationStoresListDeleted request. The method always
// closes the http.Response Body.
func (c DeletedConfigurationStoresClient) responderForConfigurationStoresListDeleted(resp *http.Response) (result ConfigurationStoresListDeletedOperationResponse, err error) {
	type page struct {
		Values   []DeletedConfigurationStore `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ConfigurationStoresListDeletedOperationResponse, err error) {
			req, err := c.preparerForConfigurationStoresListDeletedWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "deletedconfigurationstores.DeletedConfigurationStoresClient", "ConfigurationStoresListDeleted", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "deletedconfigurationstores.DeletedConfigurationStoresClient", "ConfigurationStoresListDeleted", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForConfigurationStoresListDeleted(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "deletedconfigurationstores.DeletedConfigurationStoresClient", "ConfigurationStoresListDeleted", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ConfigurationStoresListDeletedComplete retrieves all of the results into a single object
func (c DeletedConfigurationStoresClient) ConfigurationStoresListDeletedComplete(ctx context.Context, id commonids.SubscriptionId) (ConfigurationStoresListDeletedCompleteResult, error) {
	return c.ConfigurationStoresListDeletedCompleteMatchingPredicate(ctx, id, DeletedConfigurationStoreOperationPredicate{})
}

// ConfigurationStoresListDeletedCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DeletedConfigurationStoresClient) ConfigurationStoresListDeletedCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate DeletedConfigurationStoreOperationPredicate) (resp ConfigurationStoresListDeletedCompleteResult, err error) {
	items := make([]DeletedConfigurationStore, 0)

	page, err := c.ConfigurationStoresListDeleted(ctx, id)
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

	out := ConfigurationStoresListDeletedCompleteResult{
		Items: items,
	}
	return out, nil
}
