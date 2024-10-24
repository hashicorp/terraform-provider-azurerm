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

type ListAppServicesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AppServiceInfo

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListAppServicesOperationResponse, error)
}

type ListAppServicesCompleteResult struct {
	Items []AppServiceInfo
}

func (r ListAppServicesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListAppServicesOperationResponse) LoadMore(ctx context.Context) (resp ListAppServicesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListAppServices ...
func (c MonitorsClient) ListAppServices(ctx context.Context, id MonitorId) (resp ListAppServicesOperationResponse, err error) {
	req, err := c.preparerForListAppServices(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListAppServices", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListAppServices", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListAppServices(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListAppServices", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListAppServices prepares the ListAppServices request.
func (c MonitorsClient) preparerForListAppServices(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listAppServices", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListAppServicesWithNextLink prepares the ListAppServices request with the given nextLink token.
func (c MonitorsClient) preparerForListAppServicesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListAppServices handles the response to the ListAppServices request. The method always
// closes the http.Response Body.
func (c MonitorsClient) responderForListAppServices(resp *http.Response) (result ListAppServicesOperationResponse, err error) {
	type page struct {
		Values   []AppServiceInfo `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListAppServicesOperationResponse, err error) {
			req, err := c.preparerForListAppServicesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListAppServices", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListAppServices", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListAppServices(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListAppServices", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListAppServicesComplete retrieves all of the results into a single object
func (c MonitorsClient) ListAppServicesComplete(ctx context.Context, id MonitorId) (ListAppServicesCompleteResult, error) {
	return c.ListAppServicesCompleteMatchingPredicate(ctx, id, AppServiceInfoOperationPredicate{})
}

// ListAppServicesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c MonitorsClient) ListAppServicesCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate AppServiceInfoOperationPredicate) (resp ListAppServicesCompleteResult, err error) {
	items := make([]AppServiceInfo, 0)

	page, err := c.ListAppServices(ctx, id)
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

	out := ListAppServicesCompleteResult{
		Items: items,
	}
	return out, nil
}
