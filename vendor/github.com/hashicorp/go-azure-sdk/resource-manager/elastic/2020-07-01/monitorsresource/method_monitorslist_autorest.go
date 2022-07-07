package monitorsresource

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ElasticMonitorResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (MonitorsListOperationResponse, error)
}

type MonitorsListCompleteResult struct {
	Items []ElasticMonitorResource
}

func (r MonitorsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r MonitorsListOperationResponse) LoadMore(ctx context.Context) (resp MonitorsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// MonitorsList ...
func (c MonitorsResourceClient) MonitorsList(ctx context.Context, id commonids.SubscriptionId) (resp MonitorsListOperationResponse, err error) {
	req, err := c.preparerForMonitorsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForMonitorsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// MonitorsListComplete retrieves all of the results into a single object
func (c MonitorsResourceClient) MonitorsListComplete(ctx context.Context, id commonids.SubscriptionId) (MonitorsListCompleteResult, error) {
	return c.MonitorsListCompleteMatchingPredicate(ctx, id, ElasticMonitorResourceOperationPredicate{})
}

// MonitorsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c MonitorsResourceClient) MonitorsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ElasticMonitorResourceOperationPredicate) (resp MonitorsListCompleteResult, err error) {
	items := make([]ElasticMonitorResource, 0)

	page, err := c.MonitorsList(ctx, id)
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

	out := MonitorsListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForMonitorsList prepares the MonitorsList request.
func (c MonitorsResourceClient) preparerForMonitorsList(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Elastic/monitors", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForMonitorsListWithNextLink prepares the MonitorsList request with the given nextLink token.
func (c MonitorsResourceClient) preparerForMonitorsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForMonitorsList handles the response to the MonitorsList request. The method always
// closes the http.Response Body.
func (c MonitorsResourceClient) responderForMonitorsList(resp *http.Response) (result MonitorsListOperationResponse, err error) {
	type page struct {
		Values   []ElasticMonitorResource `json:"value"`
		NextLink *string                  `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result MonitorsListOperationResponse, err error) {
			req, err := c.preparerForMonitorsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForMonitorsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
