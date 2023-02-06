package linkedresources

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

type MonitorsListLinkedResourcesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]LinkedResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (MonitorsListLinkedResourcesOperationResponse, error)
}

type MonitorsListLinkedResourcesCompleteResult struct {
	Items []LinkedResource
}

func (r MonitorsListLinkedResourcesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r MonitorsListLinkedResourcesOperationResponse) LoadMore(ctx context.Context) (resp MonitorsListLinkedResourcesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// MonitorsListLinkedResources ...
func (c LinkedResourcesClient) MonitorsListLinkedResources(ctx context.Context, id MonitorId) (resp MonitorsListLinkedResourcesOperationResponse, err error) {
	req, err := c.preparerForMonitorsListLinkedResources(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "linkedresources.LinkedResourcesClient", "MonitorsListLinkedResources", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "linkedresources.LinkedResourcesClient", "MonitorsListLinkedResources", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForMonitorsListLinkedResources(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "linkedresources.LinkedResourcesClient", "MonitorsListLinkedResources", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForMonitorsListLinkedResources prepares the MonitorsListLinkedResources request.
func (c LinkedResourcesClient) preparerForMonitorsListLinkedResources(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listLinkedResources", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForMonitorsListLinkedResourcesWithNextLink prepares the MonitorsListLinkedResources request with the given nextLink token.
func (c LinkedResourcesClient) preparerForMonitorsListLinkedResourcesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForMonitorsListLinkedResources handles the response to the MonitorsListLinkedResources request. The method always
// closes the http.Response Body.
func (c LinkedResourcesClient) responderForMonitorsListLinkedResources(resp *http.Response) (result MonitorsListLinkedResourcesOperationResponse, err error) {
	type page struct {
		Values   []LinkedResource `json:"value"`
		NextLink *string          `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result MonitorsListLinkedResourcesOperationResponse, err error) {
			req, err := c.preparerForMonitorsListLinkedResourcesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "linkedresources.LinkedResourcesClient", "MonitorsListLinkedResources", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "linkedresources.LinkedResourcesClient", "MonitorsListLinkedResources", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForMonitorsListLinkedResources(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "linkedresources.LinkedResourcesClient", "MonitorsListLinkedResources", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// MonitorsListLinkedResourcesComplete retrieves all of the results into a single object
func (c LinkedResourcesClient) MonitorsListLinkedResourcesComplete(ctx context.Context, id MonitorId) (MonitorsListLinkedResourcesCompleteResult, error) {
	return c.MonitorsListLinkedResourcesCompleteMatchingPredicate(ctx, id, LinkedResourceOperationPredicate{})
}

// MonitorsListLinkedResourcesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c LinkedResourcesClient) MonitorsListLinkedResourcesCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate LinkedResourceOperationPredicate) (resp MonitorsListLinkedResourcesCompleteResult, err error) {
	items := make([]LinkedResource, 0)

	page, err := c.MonitorsListLinkedResources(ctx, id)
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

	out := MonitorsListLinkedResourcesCompleteResult{
		Items: items,
	}
	return out, nil
}
