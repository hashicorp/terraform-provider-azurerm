package labs

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

type ListVhdsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]LabVhd

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListVhdsOperationResponse, error)
}

type ListVhdsCompleteResult struct {
	Items []LabVhd
}

func (r ListVhdsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListVhdsOperationResponse) LoadMore(ctx context.Context) (resp ListVhdsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListVhds ...
func (c LabsClient) ListVhds(ctx context.Context, id LabId) (resp ListVhdsOperationResponse, err error) {
	req, err := c.preparerForListVhds(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "labs.LabsClient", "ListVhds", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "labs.LabsClient", "ListVhds", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListVhds(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "labs.LabsClient", "ListVhds", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListVhds prepares the ListVhds request.
func (c LabsClient) preparerForListVhds(ctx context.Context, id LabId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listVhds", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListVhdsWithNextLink prepares the ListVhds request with the given nextLink token.
func (c LabsClient) preparerForListVhdsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListVhds handles the response to the ListVhds request. The method always
// closes the http.Response Body.
func (c LabsClient) responderForListVhds(resp *http.Response) (result ListVhdsOperationResponse, err error) {
	type page struct {
		Values   []LabVhd `json:"value"`
		NextLink *string  `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListVhdsOperationResponse, err error) {
			req, err := c.preparerForListVhdsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "labs.LabsClient", "ListVhds", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "labs.LabsClient", "ListVhds", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListVhds(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "labs.LabsClient", "ListVhds", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListVhdsComplete retrieves all of the results into a single object
func (c LabsClient) ListVhdsComplete(ctx context.Context, id LabId) (ListVhdsCompleteResult, error) {
	return c.ListVhdsCompleteMatchingPredicate(ctx, id, LabVhdOperationPredicate{})
}

// ListVhdsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c LabsClient) ListVhdsCompleteMatchingPredicate(ctx context.Context, id LabId, predicate LabVhdOperationPredicate) (resp ListVhdsCompleteResult, err error) {
	items := make([]LabVhd, 0)

	page, err := c.ListVhds(ctx, id)
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

	out := ListVhdsCompleteResult{
		Items: items,
	}
	return out, nil
}
