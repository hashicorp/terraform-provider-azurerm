package containerinstance

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

type LocationListCachedImagesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]CachedImages

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (LocationListCachedImagesOperationResponse, error)
}

type LocationListCachedImagesCompleteResult struct {
	Items []CachedImages
}

func (r LocationListCachedImagesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r LocationListCachedImagesOperationResponse) LoadMore(ctx context.Context) (resp LocationListCachedImagesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// LocationListCachedImages ...
func (c ContainerInstanceClient) LocationListCachedImages(ctx context.Context, id LocationId) (resp LocationListCachedImagesOperationResponse, err error) {
	req, err := c.preparerForLocationListCachedImages(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListCachedImages", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListCachedImages", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForLocationListCachedImages(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListCachedImages", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForLocationListCachedImages prepares the LocationListCachedImages request.
func (c ContainerInstanceClient) preparerForLocationListCachedImages(ctx context.Context, id LocationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/cachedImages", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForLocationListCachedImagesWithNextLink prepares the LocationListCachedImages request with the given nextLink token.
func (c ContainerInstanceClient) preparerForLocationListCachedImagesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForLocationListCachedImages handles the response to the LocationListCachedImages request. The method always
// closes the http.Response Body.
func (c ContainerInstanceClient) responderForLocationListCachedImages(resp *http.Response) (result LocationListCachedImagesOperationResponse, err error) {
	type page struct {
		Values   []CachedImages `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result LocationListCachedImagesOperationResponse, err error) {
			req, err := c.preparerForLocationListCachedImagesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListCachedImages", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListCachedImages", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForLocationListCachedImages(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListCachedImages", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// LocationListCachedImagesComplete retrieves all of the results into a single object
func (c ContainerInstanceClient) LocationListCachedImagesComplete(ctx context.Context, id LocationId) (LocationListCachedImagesCompleteResult, error) {
	return c.LocationListCachedImagesCompleteMatchingPredicate(ctx, id, CachedImagesOperationPredicate{})
}

// LocationListCachedImagesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ContainerInstanceClient) LocationListCachedImagesCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate CachedImagesOperationPredicate) (resp LocationListCachedImagesCompleteResult, err error) {
	items := make([]CachedImages, 0)

	page, err := c.LocationListCachedImages(ctx, id)
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

	out := LocationListCachedImagesCompleteResult{
		Items: items,
	}
	return out, nil
}
