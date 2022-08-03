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

type QueryPacksListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]LogAnalyticsQueryPack

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (QueryPacksListOperationResponse, error)
}

type QueryPacksListCompleteResult struct {
	Items []LogAnalyticsQueryPack
}

func (r QueryPacksListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r QueryPacksListOperationResponse) LoadMore(ctx context.Context) (resp QueryPacksListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// QueryPacksList ...
func (c OperationalInsightsClient) QueryPacksList(ctx context.Context, id commonids.SubscriptionId) (resp QueryPacksListOperationResponse, err error) {
	req, err := c.preparerForQueryPacksList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForQueryPacksList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// QueryPacksListComplete retrieves all of the results into a single object
func (c OperationalInsightsClient) QueryPacksListComplete(ctx context.Context, id commonids.SubscriptionId) (QueryPacksListCompleteResult, error) {
	return c.QueryPacksListCompleteMatchingPredicate(ctx, id, LogAnalyticsQueryPackOperationPredicate{})
}

// QueryPacksListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c OperationalInsightsClient) QueryPacksListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate LogAnalyticsQueryPackOperationPredicate) (resp QueryPacksListCompleteResult, err error) {
	items := make([]LogAnalyticsQueryPack, 0)

	page, err := c.QueryPacksList(ctx, id)
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

	out := QueryPacksListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForQueryPacksList prepares the QueryPacksList request.
func (c OperationalInsightsClient) preparerForQueryPacksList(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
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

// preparerForQueryPacksListWithNextLink prepares the QueryPacksList request with the given nextLink token.
func (c OperationalInsightsClient) preparerForQueryPacksListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForQueryPacksList handles the response to the QueryPacksList request. The method always
// closes the http.Response Body.
func (c OperationalInsightsClient) responderForQueryPacksList(resp *http.Response) (result QueryPacksListOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result QueryPacksListOperationResponse, err error) {
			req, err := c.preparerForQueryPacksListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForQueryPacksList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueryPacksList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
