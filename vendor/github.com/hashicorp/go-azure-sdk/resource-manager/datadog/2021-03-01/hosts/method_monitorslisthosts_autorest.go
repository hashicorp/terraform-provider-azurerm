package hosts

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

type MonitorsListHostsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DatadogHost

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (MonitorsListHostsOperationResponse, error)
}

type MonitorsListHostsCompleteResult struct {
	Items []DatadogHost
}

func (r MonitorsListHostsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r MonitorsListHostsOperationResponse) LoadMore(ctx context.Context) (resp MonitorsListHostsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// MonitorsListHosts ...
func (c HostsClient) MonitorsListHosts(ctx context.Context, id MonitorId) (resp MonitorsListHostsOperationResponse, err error) {
	req, err := c.preparerForMonitorsListHosts(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hosts.HostsClient", "MonitorsListHosts", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "hosts.HostsClient", "MonitorsListHosts", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForMonitorsListHosts(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hosts.HostsClient", "MonitorsListHosts", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForMonitorsListHosts prepares the MonitorsListHosts request.
func (c HostsClient) preparerForMonitorsListHosts(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listHosts", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForMonitorsListHostsWithNextLink prepares the MonitorsListHosts request with the given nextLink token.
func (c HostsClient) preparerForMonitorsListHostsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForMonitorsListHosts handles the response to the MonitorsListHosts request. The method always
// closes the http.Response Body.
func (c HostsClient) responderForMonitorsListHosts(resp *http.Response) (result MonitorsListHostsOperationResponse, err error) {
	type page struct {
		Values   []DatadogHost `json:"value"`
		NextLink *string       `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result MonitorsListHostsOperationResponse, err error) {
			req, err := c.preparerForMonitorsListHostsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "hosts.HostsClient", "MonitorsListHosts", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "hosts.HostsClient", "MonitorsListHosts", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForMonitorsListHosts(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "hosts.HostsClient", "MonitorsListHosts", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// MonitorsListHostsComplete retrieves all of the results into a single object
func (c HostsClient) MonitorsListHostsComplete(ctx context.Context, id MonitorId) (MonitorsListHostsCompleteResult, error) {
	return c.MonitorsListHostsCompleteMatchingPredicate(ctx, id, DatadogHostOperationPredicate{})
}

// MonitorsListHostsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c HostsClient) MonitorsListHostsCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate DatadogHostOperationPredicate) (resp MonitorsListHostsCompleteResult, err error) {
	items := make([]DatadogHost, 0)

	page, err := c.MonitorsListHosts(ctx, id)
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

	out := MonitorsListHostsCompleteResult{
		Items: items,
	}
	return out, nil
}
