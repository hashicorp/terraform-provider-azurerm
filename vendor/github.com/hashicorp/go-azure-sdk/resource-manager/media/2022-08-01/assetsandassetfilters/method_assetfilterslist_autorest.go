package assetsandassetfilters

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

type AssetFiltersListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AssetFilter

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AssetFiltersListOperationResponse, error)
}

type AssetFiltersListCompleteResult struct {
	Items []AssetFilter
}

func (r AssetFiltersListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AssetFiltersListOperationResponse) LoadMore(ctx context.Context) (resp AssetFiltersListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// AssetFiltersList ...
func (c AssetsAndAssetFiltersClient) AssetFiltersList(ctx context.Context, id AssetId) (resp AssetFiltersListOperationResponse, err error) {
	req, err := c.preparerForAssetFiltersList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAssetFiltersList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForAssetFiltersList prepares the AssetFiltersList request.
func (c AssetsAndAssetFiltersClient) preparerForAssetFiltersList(ctx context.Context, id AssetId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/assetFilters", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForAssetFiltersListWithNextLink prepares the AssetFiltersList request with the given nextLink token.
func (c AssetsAndAssetFiltersClient) preparerForAssetFiltersListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAssetFiltersList handles the response to the AssetFiltersList request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForAssetFiltersList(resp *http.Response) (result AssetFiltersListOperationResponse, err error) {
	type page struct {
		Values   []AssetFilter `json:"value"`
		NextLink *string       `json:"@odata.nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AssetFiltersListOperationResponse, err error) {
			req, err := c.preparerForAssetFiltersListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAssetFiltersList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetFiltersList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// AssetFiltersListComplete retrieves all of the results into a single object
func (c AssetsAndAssetFiltersClient) AssetFiltersListComplete(ctx context.Context, id AssetId) (AssetFiltersListCompleteResult, error) {
	return c.AssetFiltersListCompleteMatchingPredicate(ctx, id, AssetFilterOperationPredicate{})
}

// AssetFiltersListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AssetsAndAssetFiltersClient) AssetFiltersListCompleteMatchingPredicate(ctx context.Context, id AssetId, predicate AssetFilterOperationPredicate) (resp AssetFiltersListCompleteResult, err error) {
	items := make([]AssetFilter, 0)

	page, err := c.AssetFiltersList(ctx, id)
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

	out := AssetFiltersListCompleteResult{
		Items: items,
	}
	return out, nil
}
