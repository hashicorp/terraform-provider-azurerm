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

type AssetsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Asset

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AssetsListOperationResponse, error)
}

type AssetsListCompleteResult struct {
	Items []Asset
}

func (r AssetsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AssetsListOperationResponse) LoadMore(ctx context.Context) (resp AssetsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type AssetsListOperationOptions struct {
	Filter  *string
	Orderby *string
	Top     *int64
}

func DefaultAssetsListOperationOptions() AssetsListOperationOptions {
	return AssetsListOperationOptions{}
}

func (o AssetsListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o AssetsListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Orderby != nil {
		out["$orderby"] = *o.Orderby
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// AssetsList ...
func (c AssetsAndAssetFiltersClient) AssetsList(ctx context.Context, id MediaServiceId, options AssetsListOperationOptions) (resp AssetsListOperationResponse, err error) {
	req, err := c.preparerForAssetsList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAssetsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForAssetsList prepares the AssetsList request.
func (c AssetsAndAssetFiltersClient) preparerForAssetsList(ctx context.Context, id MediaServiceId, options AssetsListOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/assets", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForAssetsListWithNextLink prepares the AssetsList request with the given nextLink token.
func (c AssetsAndAssetFiltersClient) preparerForAssetsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAssetsList handles the response to the AssetsList request. The method always
// closes the http.Response Body.
func (c AssetsAndAssetFiltersClient) responderForAssetsList(resp *http.Response) (result AssetsListOperationResponse, err error) {
	type page struct {
		Values   []Asset `json:"value"`
		NextLink *string `json:"@odata.nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AssetsListOperationResponse, err error) {
			req, err := c.preparerForAssetsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAssetsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "assetsandassetfilters.AssetsAndAssetFiltersClient", "AssetsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// AssetsListComplete retrieves all of the results into a single object
func (c AssetsAndAssetFiltersClient) AssetsListComplete(ctx context.Context, id MediaServiceId, options AssetsListOperationOptions) (AssetsListCompleteResult, error) {
	return c.AssetsListCompleteMatchingPredicate(ctx, id, options, AssetOperationPredicate{})
}

// AssetsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AssetsAndAssetFiltersClient) AssetsListCompleteMatchingPredicate(ctx context.Context, id MediaServiceId, options AssetsListOperationOptions, predicate AssetOperationPredicate) (resp AssetsListCompleteResult, err error) {
	items := make([]Asset, 0)

	page, err := c.AssetsList(ctx, id, options)
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

	out := AssetsListCompleteResult{
		Items: items,
	}
	return out, nil
}
