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

type LocationListCapabilitiesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Capabilities

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (LocationListCapabilitiesOperationResponse, error)
}

type LocationListCapabilitiesCompleteResult struct {
	Items []Capabilities
}

func (r LocationListCapabilitiesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r LocationListCapabilitiesOperationResponse) LoadMore(ctx context.Context) (resp LocationListCapabilitiesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// LocationListCapabilities ...
func (c ContainerInstanceClient) LocationListCapabilities(ctx context.Context, id LocationId) (resp LocationListCapabilitiesOperationResponse, err error) {
	req, err := c.preparerForLocationListCapabilities(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListCapabilities", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListCapabilities", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForLocationListCapabilities(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListCapabilities", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForLocationListCapabilities prepares the LocationListCapabilities request.
func (c ContainerInstanceClient) preparerForLocationListCapabilities(ctx context.Context, id LocationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/capabilities", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForLocationListCapabilitiesWithNextLink prepares the LocationListCapabilities request with the given nextLink token.
func (c ContainerInstanceClient) preparerForLocationListCapabilitiesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForLocationListCapabilities handles the response to the LocationListCapabilities request. The method always
// closes the http.Response Body.
func (c ContainerInstanceClient) responderForLocationListCapabilities(resp *http.Response) (result LocationListCapabilitiesOperationResponse, err error) {
	type page struct {
		Values   []Capabilities `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result LocationListCapabilitiesOperationResponse, err error) {
			req, err := c.preparerForLocationListCapabilitiesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListCapabilities", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListCapabilities", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForLocationListCapabilities(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListCapabilities", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// LocationListCapabilitiesComplete retrieves all of the results into a single object
func (c ContainerInstanceClient) LocationListCapabilitiesComplete(ctx context.Context, id LocationId) (LocationListCapabilitiesCompleteResult, error) {
	return c.LocationListCapabilitiesCompleteMatchingPredicate(ctx, id, CapabilitiesOperationPredicate{})
}

// LocationListCapabilitiesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ContainerInstanceClient) LocationListCapabilitiesCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate CapabilitiesOperationPredicate) (resp LocationListCapabilitiesCompleteResult, err error) {
	items := make([]Capabilities, 0)

	page, err := c.LocationListCapabilities(ctx, id)
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

	out := LocationListCapabilitiesCompleteResult{
		Items: items,
	}
	return out, nil
}
