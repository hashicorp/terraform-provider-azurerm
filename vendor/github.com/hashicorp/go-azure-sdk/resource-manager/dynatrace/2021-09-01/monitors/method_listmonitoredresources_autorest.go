package monitors

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

type ListMonitoredResourcesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]MonitoredResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListMonitoredResourcesOperationResponse, error)
}

type ListMonitoredResourcesCompleteResult struct {
	Items []MonitoredResource
}

func (r ListMonitoredResourcesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListMonitoredResourcesOperationResponse) LoadMore(ctx context.Context) (resp ListMonitoredResourcesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListMonitoredResources ...
func (c MonitorsClient) ListMonitoredResources(ctx context.Context, id MonitorId) (resp ListMonitoredResourcesOperationResponse, err error) {
	req, err := c.preparerForListMonitoredResources(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListMonitoredResources", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListMonitoredResources", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListMonitoredResources(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListMonitoredResources", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListMonitoredResources prepares the ListMonitoredResources request.
func (c MonitorsClient) preparerForListMonitoredResources(ctx context.Context, id MonitorId) (*http.Request, error) {
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

// preparerForListMonitoredResourcesWithNextLink prepares the ListMonitoredResources request with the given nextLink token.
func (c MonitorsClient) preparerForListMonitoredResourcesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListMonitoredResources handles the response to the ListMonitoredResources request. The method always
// closes the http.Response Body.
func (c MonitorsClient) responderForListMonitoredResources(resp *http.Response) (result ListMonitoredResourcesOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListMonitoredResourcesOperationResponse, err error) {
			req, err := c.preparerForListMonitoredResourcesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListMonitoredResources", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListMonitoredResources", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListMonitoredResources(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListMonitoredResources", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListMonitoredResourcesComplete retrieves all of the results into a single object
func (c MonitorsClient) ListMonitoredResourcesComplete(ctx context.Context, id MonitorId) (ListMonitoredResourcesCompleteResult, error) {
	return c.ListMonitoredResourcesCompleteMatchingPredicate(ctx, id, MonitoredResourceOperationPredicate{})
}

// ListMonitoredResourcesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c MonitorsClient) ListMonitoredResourcesCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate MonitoredResourceOperationPredicate) (resp ListMonitoredResourcesCompleteResult, err error) {
	items := make([]MonitoredResource, 0)

	page, err := c.ListMonitoredResources(ctx, id)
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

	out := ListMonitoredResourcesCompleteResult{
		Items: items,
	}
	return out, nil
}
