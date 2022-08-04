package operationalinsights

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

type QueryPacksListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]LogAnalyticsQueryPack

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (QueryPacksListByResourceGroupOperationResponse, error)
}

type QueryPacksListByResourceGroupCompleteResult struct {
	Items []LogAnalyticsQueryPack
}

func (r QueryPacksListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r QueryPacksListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp QueryPacksListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// QueryPacksListByResourceGroup ...
func (c OperationalInsightsClient) QueryPacksListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp QueryPacksListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForQueryPacksListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForQueryPacksListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// QueryPacksListByResourceGroupComplete retrieves all of the results into a single object
func (c OperationalInsightsClient) QueryPacksListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (QueryPacksListByResourceGroupCompleteResult, error) {
	return c.QueryPacksListByResourceGroupCompleteMatchingPredicate(ctx, id, LogAnalyticsQueryPackOperationPredicate{})
}

// QueryPacksListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c OperationalInsightsClient) QueryPacksListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate LogAnalyticsQueryPackOperationPredicate) (resp QueryPacksListByResourceGroupCompleteResult, err error) {
	items := make([]LogAnalyticsQueryPack, 0)

	page, err := c.QueryPacksListByResourceGroup(ctx, id)
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

	out := QueryPacksListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForQueryPacksListByResourceGroup prepares the QueryPacksListByResourceGroup request.
func (c OperationalInsightsClient) preparerForQueryPacksListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.OperationalInsights/queryPacks", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForQueryPacksListByResourceGroupWithNextLink prepares the QueryPacksListByResourceGroup request with the given nextLink token.
func (c OperationalInsightsClient) preparerForQueryPacksListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForQueryPacksListByResourceGroup handles the response to the QueryPacksListByResourceGroup request. The method always
// closes the http.Response Body.
func (c OperationalInsightsClient) responderForQueryPacksListByResourceGroup(resp *http.Response) (result QueryPacksListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []LogAnalyticsQueryPack `json:"value"`
		NextLink *string                 `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result QueryPacksListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForQueryPacksListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForQueryPacksListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
