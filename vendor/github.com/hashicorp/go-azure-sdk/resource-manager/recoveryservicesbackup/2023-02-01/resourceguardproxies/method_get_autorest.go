package resourceguardproxies

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

type GetOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ResourceGuardProxyBaseResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (GetOperationResponse, error)
}

type GetCompleteResult struct {
	Items []ResourceGuardProxyBaseResource
}

func (r GetOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r GetOperationResponse) LoadMore(ctx context.Context) (resp GetOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// Get ...
func (c ResourceGuardProxiesClient) Get(ctx context.Context, id VaultId) (resp GetOperationResponse, err error) {
	req, err := c.preparerForGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguardproxies.ResourceGuardProxiesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguardproxies.ResourceGuardProxiesClient", "Get", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForGet(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguardproxies.ResourceGuardProxiesClient", "Get", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForGet prepares the Get request.
func (c ResourceGuardProxiesClient) preparerForGet(ctx context.Context, id VaultId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/backupResourceGuardProxies", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForGetWithNextLink prepares the Get request with the given nextLink token.
func (c ResourceGuardProxiesClient) preparerForGetWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForGet handles the response to the Get request. The method always
// closes the http.Response Body.
func (c ResourceGuardProxiesClient) responderForGet(resp *http.Response) (result GetOperationResponse, err error) {
	type page struct {
		Values   []ResourceGuardProxyBaseResource `json:"value"`
		NextLink *string                          `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result GetOperationResponse, err error) {
			req, err := c.preparerForGetWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguardproxies.ResourceGuardProxiesClient", "Get", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguardproxies.ResourceGuardProxiesClient", "Get", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForGet(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguardproxies.ResourceGuardProxiesClient", "Get", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// GetComplete retrieves all of the results into a single object
func (c ResourceGuardProxiesClient) GetComplete(ctx context.Context, id VaultId) (GetCompleteResult, error) {
	return c.GetCompleteMatchingPredicate(ctx, id, ResourceGuardProxyBaseResourceOperationPredicate{})
}

// GetCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ResourceGuardProxiesClient) GetCompleteMatchingPredicate(ctx context.Context, id VaultId, predicate ResourceGuardProxyBaseResourceOperationPredicate) (resp GetCompleteResult, err error) {
	items := make([]ResourceGuardProxyBaseResource, 0)

	page, err := c.Get(ctx, id)
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

	out := GetCompleteResult{
		Items: items,
	}
	return out, nil
}
