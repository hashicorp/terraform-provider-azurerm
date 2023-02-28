package monitoredresources

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

type MonitorsListMonitoredResourcesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]MonitoredResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (MonitorsListMonitoredResourcesOperationResponse, error)
}

type MonitorsListMonitoredResourcesCompleteResult struct {
	Items []MonitoredResource
}

func (r MonitorsListMonitoredResourcesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r MonitorsListMonitoredResourcesOperationResponse) LoadMore(ctx context.Context) (resp MonitorsListMonitoredResourcesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// MonitorsListMonitoredResources ...
func (c MonitoredResourcesClient) MonitorsListMonitoredResources(ctx context.Context, id MonitorId) (resp MonitorsListMonitoredResourcesOperationResponse, err error) {
	req, err := c.preparerForMonitorsListMonitoredResources(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitoredresources.MonitoredResourcesClient", "MonitorsListMonitoredResources", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitoredresources.MonitoredResourcesClient", "MonitorsListMonitoredResources", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForMonitorsListMonitoredResources(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitoredresources.MonitoredResourcesClient", "MonitorsListMonitoredResources", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForMonitorsListMonitoredResources prepares the MonitorsListMonitoredResources request.
func (c MonitoredResourcesClient) preparerForMonitorsListMonitoredResources(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listMonitoredResources", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForMonitorsListMonitoredResourcesWithNextLink prepares the MonitorsListMonitoredResources request with the given nextLink token.
func (c MonitoredResourcesClient) preparerForMonitorsListMonitoredResourcesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForMonitorsListMonitoredResources handles the response to the MonitorsListMonitoredResources request. The method always
// closes the http.Response Body.
func (c MonitoredResourcesClient) responderForMonitorsListMonitoredResources(resp *http.Response) (result MonitorsListMonitoredResourcesOperationResponse, err error) {
	type page struct {
		Values   []MonitoredResource `json:"value"`
		NextLink *string             `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result MonitorsListMonitoredResourcesOperationResponse, err error) {
			req, err := c.preparerForMonitorsListMonitoredResourcesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitoredresources.MonitoredResourcesClient", "MonitorsListMonitoredResources", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitoredresources.MonitoredResourcesClient", "MonitorsListMonitoredResources", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForMonitorsListMonitoredResources(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitoredresources.MonitoredResourcesClient", "MonitorsListMonitoredResources", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// MonitorsListMonitoredResourcesComplete retrieves all of the results into a single object
func (c MonitoredResourcesClient) MonitorsListMonitoredResourcesComplete(ctx context.Context, id MonitorId) (MonitorsListMonitoredResourcesCompleteResult, error) {
	return c.MonitorsListMonitoredResourcesCompleteMatchingPredicate(ctx, id, MonitoredResourceOperationPredicate{})
}

// MonitorsListMonitoredResourcesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c MonitoredResourcesClient) MonitorsListMonitoredResourcesCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate MonitoredResourceOperationPredicate) (resp MonitorsListMonitoredResourcesCompleteResult, err error) {
	items := make([]MonitoredResource, 0)

	page, err := c.MonitorsListMonitoredResources(ctx, id)
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

	out := MonitorsListMonitoredResourcesCompleteResult{
		Items: items,
	}
	return out, nil
}
