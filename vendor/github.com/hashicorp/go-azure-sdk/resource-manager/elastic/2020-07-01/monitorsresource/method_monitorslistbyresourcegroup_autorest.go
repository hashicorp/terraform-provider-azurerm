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

type MonitorsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ElasticMonitorResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (MonitorsListByResourceGroupOperationResponse, error)
}

type MonitorsListByResourceGroupCompleteResult struct {
	Items []ElasticMonitorResource
}

func (r MonitorsListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r MonitorsListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp MonitorsListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// MonitorsListByResourceGroup ...
func (c MonitorsResourceClient) MonitorsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp MonitorsListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForMonitorsListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForMonitorsListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// MonitorsListByResourceGroupComplete retrieves all of the results into a single object
func (c MonitorsResourceClient) MonitorsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (MonitorsListByResourceGroupCompleteResult, error) {
	return c.MonitorsListByResourceGroupCompleteMatchingPredicate(ctx, id, ElasticMonitorResourceOperationPredicate{})
}

// MonitorsListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c MonitorsResourceClient) MonitorsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ElasticMonitorResourceOperationPredicate) (resp MonitorsListByResourceGroupCompleteResult, err error) {
	items := make([]ElasticMonitorResource, 0)

	page, err := c.MonitorsListByResourceGroup(ctx, id)
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

	out := MonitorsListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForMonitorsListByResourceGroup prepares the MonitorsListByResourceGroup request.
func (c MonitorsResourceClient) preparerForMonitorsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
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

// preparerForMonitorsListByResourceGroupWithNextLink prepares the MonitorsListByResourceGroup request with the given nextLink token.
func (c MonitorsResourceClient) preparerForMonitorsListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForMonitorsListByResourceGroup handles the response to the MonitorsListByResourceGroup request. The method always
// closes the http.Response Body.
func (c MonitorsResourceClient) responderForMonitorsListByResourceGroup(resp *http.Response) (result MonitorsListByResourceGroupOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result MonitorsListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForMonitorsListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForMonitorsListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitorsresource.MonitorsResourceClient", "MonitorsListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
