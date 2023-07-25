package frontdoors

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

type FrontendEndpointsListByFrontDoorOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]FrontendEndpoint

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (FrontendEndpointsListByFrontDoorOperationResponse, error)
}

type FrontendEndpointsListByFrontDoorCompleteResult struct {
	Items []FrontendEndpoint
}

func (r FrontendEndpointsListByFrontDoorOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r FrontendEndpointsListByFrontDoorOperationResponse) LoadMore(ctx context.Context) (resp FrontendEndpointsListByFrontDoorOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// FrontendEndpointsListByFrontDoor ...
func (c FrontDoorsClient) FrontendEndpointsListByFrontDoor(ctx context.Context, id FrontDoorId) (resp FrontendEndpointsListByFrontDoorOperationResponse, err error) {
	req, err := c.preparerForFrontendEndpointsListByFrontDoor(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsListByFrontDoor", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsListByFrontDoor", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForFrontendEndpointsListByFrontDoor(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsListByFrontDoor", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForFrontendEndpointsListByFrontDoor prepares the FrontendEndpointsListByFrontDoor request.
func (c FrontDoorsClient) preparerForFrontendEndpointsListByFrontDoor(ctx context.Context, id FrontDoorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/frontendEndpoints", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForFrontendEndpointsListByFrontDoorWithNextLink prepares the FrontendEndpointsListByFrontDoor request with the given nextLink token.
func (c FrontDoorsClient) preparerForFrontendEndpointsListByFrontDoorWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForFrontendEndpointsListByFrontDoor handles the response to the FrontendEndpointsListByFrontDoor request. The method always
// closes the http.Response Body.
func (c FrontDoorsClient) responderForFrontendEndpointsListByFrontDoor(resp *http.Response) (result FrontendEndpointsListByFrontDoorOperationResponse, err error) {
	type page struct {
		Values   []FrontendEndpoint `json:"value"`
		NextLink *string            `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result FrontendEndpointsListByFrontDoorOperationResponse, err error) {
			req, err := c.preparerForFrontendEndpointsListByFrontDoorWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsListByFrontDoor", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsListByFrontDoor", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForFrontendEndpointsListByFrontDoor(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "FrontendEndpointsListByFrontDoor", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// FrontendEndpointsListByFrontDoorComplete retrieves all of the results into a single object
func (c FrontDoorsClient) FrontendEndpointsListByFrontDoorComplete(ctx context.Context, id FrontDoorId) (FrontendEndpointsListByFrontDoorCompleteResult, error) {
	return c.FrontendEndpointsListByFrontDoorCompleteMatchingPredicate(ctx, id, FrontendEndpointOperationPredicate{})
}

// FrontendEndpointsListByFrontDoorCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c FrontDoorsClient) FrontendEndpointsListByFrontDoorCompleteMatchingPredicate(ctx context.Context, id FrontDoorId, predicate FrontendEndpointOperationPredicate) (resp FrontendEndpointsListByFrontDoorCompleteResult, err error) {
	items := make([]FrontendEndpoint, 0)

	page, err := c.FrontendEndpointsListByFrontDoor(ctx, id)
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

	out := FrontendEndpointsListByFrontDoorCompleteResult{
		Items: items,
	}
	return out, nil
}
