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

type ListHostsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]VMInfo

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListHostsOperationResponse, error)
}

type ListHostsCompleteResult struct {
	Items []VMInfo
}

func (r ListHostsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListHostsOperationResponse) LoadMore(ctx context.Context) (resp ListHostsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListHosts ...
func (c MonitorsClient) ListHosts(ctx context.Context, id MonitorId) (resp ListHostsOperationResponse, err error) {
	req, err := c.preparerForListHosts(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListHosts", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListHosts", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListHosts(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListHosts", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListHosts prepares the ListHosts request.
func (c MonitorsClient) preparerForListHosts(ctx context.Context, id MonitorId) (*http.Request, error) {
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

// preparerForListHostsWithNextLink prepares the ListHosts request with the given nextLink token.
func (c MonitorsClient) preparerForListHostsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListHosts handles the response to the ListHosts request. The method always
// closes the http.Response Body.
func (c MonitorsClient) responderForListHosts(resp *http.Response) (result ListHostsOperationResponse, err error) {
	type page struct {
		Values   []VMInfo `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListHostsOperationResponse, err error) {
			req, err := c.preparerForListHostsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListHosts", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListHosts", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListHosts(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListHosts", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListHostsComplete retrieves all of the results into a single object
func (c MonitorsClient) ListHostsComplete(ctx context.Context, id MonitorId) (ListHostsCompleteResult, error) {
	return c.ListHostsCompleteMatchingPredicate(ctx, id, VMInfoOperationPredicate{})
}

// ListHostsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c MonitorsClient) ListHostsCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate VMInfoOperationPredicate) (resp ListHostsCompleteResult, err error) {
	items := make([]VMInfo, 0)

	page, err := c.ListHosts(ctx, id)
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

	out := ListHostsCompleteResult{
		Items: items,
	}
	return out, nil
}
