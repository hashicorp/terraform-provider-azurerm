package videoanalyzer

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

type VideosListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]VideoEntity

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (VideosListOperationResponse, error)
}

type VideosListCompleteResult struct {
	Items []VideoEntity
}

func (r VideosListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r VideosListOperationResponse) LoadMore(ctx context.Context) (resp VideosListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type VideosListOperationOptions struct {
	Top *int64
}

func DefaultVideosListOperationOptions() VideosListOperationOptions {
	return VideosListOperationOptions{}
}

func (o VideosListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o VideosListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// VideosList ...
func (c VideoAnalyzerClient) VideosList(ctx context.Context, id VideoAnalyzerId, options VideosListOperationOptions) (resp VideosListOperationResponse, err error) {
	req, err := c.preparerForVideosList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForVideosList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// VideosListComplete retrieves all of the results into a single object
func (c VideoAnalyzerClient) VideosListComplete(ctx context.Context, id VideoAnalyzerId, options VideosListOperationOptions) (VideosListCompleteResult, error) {
	return c.VideosListCompleteMatchingPredicate(ctx, id, options, VideoEntityOperationPredicate{})
}

// VideosListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c VideoAnalyzerClient) VideosListCompleteMatchingPredicate(ctx context.Context, id VideoAnalyzerId, options VideosListOperationOptions, predicate VideoEntityOperationPredicate) (resp VideosListCompleteResult, err error) {
	items := make([]VideoEntity, 0)

	page, err := c.VideosList(ctx, id, options)
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

	out := VideosListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForVideosList prepares the VideosList request.
func (c VideoAnalyzerClient) preparerForVideosList(ctx context.Context, id VideoAnalyzerId, options VideosListOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/videos", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForVideosListWithNextLink prepares the VideosList request with the given nextLink token.
func (c VideoAnalyzerClient) preparerForVideosListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForVideosList handles the response to the VideosList request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForVideosList(resp *http.Response) (result VideosListOperationResponse, err error) {
	type page struct {
		Values   []VideoEntity `json:"value"`
		NextLink *string       `json:"@nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result VideosListOperationResponse, err error) {
			req, err := c.preparerForVideosListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForVideosList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "VideosList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
