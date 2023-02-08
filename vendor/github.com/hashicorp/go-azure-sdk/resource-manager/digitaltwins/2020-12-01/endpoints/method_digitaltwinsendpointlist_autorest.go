package endpoints

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

type DigitalTwinsEndpointListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DigitalTwinsEndpointResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (DigitalTwinsEndpointListOperationResponse, error)
}

type DigitalTwinsEndpointListCompleteResult struct {
	Items []DigitalTwinsEndpointResource
}

func (r DigitalTwinsEndpointListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r DigitalTwinsEndpointListOperationResponse) LoadMore(ctx context.Context) (resp DigitalTwinsEndpointListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// DigitalTwinsEndpointList ...
func (c EndpointsClient) DigitalTwinsEndpointList(ctx context.Context, id DigitalTwinsInstanceId) (resp DigitalTwinsEndpointListOperationResponse, err error) {
	req, err := c.preparerForDigitalTwinsEndpointList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "DigitalTwinsEndpointList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "DigitalTwinsEndpointList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForDigitalTwinsEndpointList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "DigitalTwinsEndpointList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForDigitalTwinsEndpointList prepares the DigitalTwinsEndpointList request.
func (c EndpointsClient) preparerForDigitalTwinsEndpointList(ctx context.Context, id DigitalTwinsInstanceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/endpoints", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForDigitalTwinsEndpointListWithNextLink prepares the DigitalTwinsEndpointList request with the given nextLink token.
func (c EndpointsClient) preparerForDigitalTwinsEndpointListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForDigitalTwinsEndpointList handles the response to the DigitalTwinsEndpointList request. The method always
// closes the http.Response Body.
func (c EndpointsClient) responderForDigitalTwinsEndpointList(resp *http.Response) (result DigitalTwinsEndpointListOperationResponse, err error) {
	type page struct {
		Values   []DigitalTwinsEndpointResource `json:"value"`
		NextLink *string                        `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result DigitalTwinsEndpointListOperationResponse, err error) {
			req, err := c.preparerForDigitalTwinsEndpointListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "DigitalTwinsEndpointList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "DigitalTwinsEndpointList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForDigitalTwinsEndpointList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "endpoints.EndpointsClient", "DigitalTwinsEndpointList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// DigitalTwinsEndpointListComplete retrieves all of the results into a single object
func (c EndpointsClient) DigitalTwinsEndpointListComplete(ctx context.Context, id DigitalTwinsInstanceId) (DigitalTwinsEndpointListCompleteResult, error) {
	return c.DigitalTwinsEndpointListCompleteMatchingPredicate(ctx, id, DigitalTwinsEndpointResourceOperationPredicate{})
}

// DigitalTwinsEndpointListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EndpointsClient) DigitalTwinsEndpointListCompleteMatchingPredicate(ctx context.Context, id DigitalTwinsInstanceId, predicate DigitalTwinsEndpointResourceOperationPredicate) (resp DigitalTwinsEndpointListCompleteResult, err error) {
	items := make([]DigitalTwinsEndpointResource, 0)

	page, err := c.DigitalTwinsEndpointList(ctx, id)
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

	out := DigitalTwinsEndpointListCompleteResult{
		Items: items,
	}
	return out, nil
}
