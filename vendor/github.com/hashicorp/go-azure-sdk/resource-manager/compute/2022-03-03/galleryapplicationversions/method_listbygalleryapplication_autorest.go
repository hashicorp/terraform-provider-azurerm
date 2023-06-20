package galleryapplicationversions

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

type ListByGalleryApplicationOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]GalleryApplicationVersion

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByGalleryApplicationOperationResponse, error)
}

type ListByGalleryApplicationCompleteResult struct {
	Items []GalleryApplicationVersion
}

func (r ListByGalleryApplicationOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByGalleryApplicationOperationResponse) LoadMore(ctx context.Context) (resp ListByGalleryApplicationOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByGalleryApplication ...
func (c GalleryApplicationVersionsClient) ListByGalleryApplication(ctx context.Context, id ApplicationId) (resp ListByGalleryApplicationOperationResponse, err error) {
	req, err := c.preparerForListByGalleryApplication(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "galleryapplicationversions.GalleryApplicationVersionsClient", "ListByGalleryApplication", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "galleryapplicationversions.GalleryApplicationVersionsClient", "ListByGalleryApplication", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByGalleryApplication(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "galleryapplicationversions.GalleryApplicationVersionsClient", "ListByGalleryApplication", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByGalleryApplication prepares the ListByGalleryApplication request.
func (c GalleryApplicationVersionsClient) preparerForListByGalleryApplication(ctx context.Context, id ApplicationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/versions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByGalleryApplicationWithNextLink prepares the ListByGalleryApplication request with the given nextLink token.
func (c GalleryApplicationVersionsClient) preparerForListByGalleryApplicationWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByGalleryApplication handles the response to the ListByGalleryApplication request. The method always
// closes the http.Response Body.
func (c GalleryApplicationVersionsClient) responderForListByGalleryApplication(resp *http.Response) (result ListByGalleryApplicationOperationResponse, err error) {
	type page struct {
		Values   []GalleryApplicationVersion `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByGalleryApplicationOperationResponse, err error) {
			req, err := c.preparerForListByGalleryApplicationWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "galleryapplicationversions.GalleryApplicationVersionsClient", "ListByGalleryApplication", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "galleryapplicationversions.GalleryApplicationVersionsClient", "ListByGalleryApplication", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByGalleryApplication(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "galleryapplicationversions.GalleryApplicationVersionsClient", "ListByGalleryApplication", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByGalleryApplicationComplete retrieves all of the results into a single object
func (c GalleryApplicationVersionsClient) ListByGalleryApplicationComplete(ctx context.Context, id ApplicationId) (ListByGalleryApplicationCompleteResult, error) {
	return c.ListByGalleryApplicationCompleteMatchingPredicate(ctx, id, GalleryApplicationVersionOperationPredicate{})
}

// ListByGalleryApplicationCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c GalleryApplicationVersionsClient) ListByGalleryApplicationCompleteMatchingPredicate(ctx context.Context, id ApplicationId, predicate GalleryApplicationVersionOperationPredicate) (resp ListByGalleryApplicationCompleteResult, err error) {
	items := make([]GalleryApplicationVersion, 0)

	page, err := c.ListByGalleryApplication(ctx, id)
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

	out := ListByGalleryApplicationCompleteResult{
		Items: items,
	}
	return out, nil
}
