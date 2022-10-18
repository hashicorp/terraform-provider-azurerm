package grafanaresource

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

type GrafanaListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ManagedGrafana

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (GrafanaListByResourceGroupOperationResponse, error)
}

type GrafanaListByResourceGroupCompleteResult struct {
	Items []ManagedGrafana
}

func (r GrafanaListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r GrafanaListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp GrafanaListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// GrafanaListByResourceGroup ...
func (c GrafanaResourceClient) GrafanaListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp GrafanaListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForGrafanaListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForGrafanaListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForGrafanaListByResourceGroup prepares the GrafanaListByResourceGroup request.
func (c GrafanaResourceClient) preparerForGrafanaListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Dashboard/grafana", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForGrafanaListByResourceGroupWithNextLink prepares the GrafanaListByResourceGroup request with the given nextLink token.
func (c GrafanaResourceClient) preparerForGrafanaListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForGrafanaListByResourceGroup handles the response to the GrafanaListByResourceGroup request. The method always
// closes the http.Response Body.
func (c GrafanaResourceClient) responderForGrafanaListByResourceGroup(resp *http.Response) (result GrafanaListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []ManagedGrafana `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result GrafanaListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForGrafanaListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForGrafanaListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "grafanaresource.GrafanaResourceClient", "GrafanaListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// GrafanaListByResourceGroupComplete retrieves all of the results into a single object
func (c GrafanaResourceClient) GrafanaListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (GrafanaListByResourceGroupCompleteResult, error) {
	return c.GrafanaListByResourceGroupCompleteMatchingPredicate(ctx, id, ManagedGrafanaOperationPredicate{})
}

// GrafanaListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c GrafanaResourceClient) GrafanaListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ManagedGrafanaOperationPredicate) (resp GrafanaListByResourceGroupCompleteResult, err error) {
	items := make([]ManagedGrafana, 0)

	page, err := c.GrafanaListByResourceGroup(ctx, id)
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

	out := GrafanaListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
