package encodings

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

type TransformsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Transform

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (TransformsListOperationResponse, error)
}

type TransformsListCompleteResult struct {
	Items []Transform
}

func (r TransformsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r TransformsListOperationResponse) LoadMore(ctx context.Context) (resp TransformsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type TransformsListOperationOptions struct {
	Filter  *string
	Orderby *string
}

func DefaultTransformsListOperationOptions() TransformsListOperationOptions {
	return TransformsListOperationOptions{}
}

func (o TransformsListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o TransformsListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Orderby != nil {
		out["$orderby"] = *o.Orderby
	}

	return out
}

// TransformsList ...
func (c EncodingsClient) TransformsList(ctx context.Context, id MediaServiceId, options TransformsListOperationOptions) (resp TransformsListOperationResponse, err error) {
	req, err := c.preparerForTransformsList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForTransformsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForTransformsList prepares the TransformsList request.
func (c EncodingsClient) preparerForTransformsList(ctx context.Context, id MediaServiceId, options TransformsListOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/transforms", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForTransformsListWithNextLink prepares the TransformsList request with the given nextLink token.
func (c EncodingsClient) preparerForTransformsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForTransformsList handles the response to the TransformsList request. The method always
// closes the http.Response Body.
func (c EncodingsClient) responderForTransformsList(resp *http.Response) (result TransformsListOperationResponse, err error) {
	type page struct {
		Values   []Transform `json:"value"`
		NextLink *string     `json:"@odata.nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result TransformsListOperationResponse, err error) {
			req, err := c.preparerForTransformsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForTransformsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "encodings.EncodingsClient", "TransformsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// TransformsListComplete retrieves all of the results into a single object
func (c EncodingsClient) TransformsListComplete(ctx context.Context, id MediaServiceId, options TransformsListOperationOptions) (TransformsListCompleteResult, error) {
	return c.TransformsListCompleteMatchingPredicate(ctx, id, options, TransformOperationPredicate{})
}

// TransformsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EncodingsClient) TransformsListCompleteMatchingPredicate(ctx context.Context, id MediaServiceId, options TransformsListOperationOptions, predicate TransformOperationPredicate) (resp TransformsListCompleteResult, err error) {
	items := make([]Transform, 0)

	page, err := c.TransformsList(ctx, id, options)
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

	out := TransformsListCompleteResult{
		Items: items,
	}
	return out, nil
}
