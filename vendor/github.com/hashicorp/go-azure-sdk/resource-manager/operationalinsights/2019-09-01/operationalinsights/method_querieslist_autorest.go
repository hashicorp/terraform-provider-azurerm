package operationalinsights

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

type QueriesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]LogAnalyticsQueryPackQuery

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (QueriesListOperationResponse, error)
}

type QueriesListCompleteResult struct {
	Items []LogAnalyticsQueryPackQuery
}

func (r QueriesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r QueriesListOperationResponse) LoadMore(ctx context.Context) (resp QueriesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type QueriesListOperationOptions struct {
	IncludeBody *bool
	Top         *int64
}

func DefaultQueriesListOperationOptions() QueriesListOperationOptions {
	return QueriesListOperationOptions{}
}

func (o QueriesListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o QueriesListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IncludeBody != nil {
		out["includeBody"] = *o.IncludeBody
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// QueriesList ...
func (c OperationalInsightsClient) QueriesList(ctx context.Context, id QueryPackId, options QueriesListOperationOptions) (resp QueriesListOperationResponse, err error) {
	req, err := c.preparerForQueriesList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForQueriesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForQueriesList prepares the QueriesList request.
func (c OperationalInsightsClient) preparerForQueriesList(ctx context.Context, id QueryPackId, options QueriesListOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/queries", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForQueriesListWithNextLink prepares the QueriesList request with the given nextLink token.
func (c OperationalInsightsClient) preparerForQueriesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForQueriesList handles the response to the QueriesList request. The method always
// closes the http.Response Body.
func (c OperationalInsightsClient) responderForQueriesList(resp *http.Response) (result QueriesListOperationResponse, err error) {
	type page struct {
		Values   []LogAnalyticsQueryPackQuery `json:"value"`
		NextLink *string                      `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result QueriesListOperationResponse, err error) {
			req, err := c.preparerForQueriesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForQueriesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// QueriesListComplete retrieves all of the results into a single object
func (c OperationalInsightsClient) QueriesListComplete(ctx context.Context, id QueryPackId, options QueriesListOperationOptions) (QueriesListCompleteResult, error) {
	return c.QueriesListCompleteMatchingPredicate(ctx, id, options, LogAnalyticsQueryPackQueryOperationPredicate{})
}

// QueriesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c OperationalInsightsClient) QueriesListCompleteMatchingPredicate(ctx context.Context, id QueryPackId, options QueriesListOperationOptions, predicate LogAnalyticsQueryPackQueryOperationPredicate) (resp QueriesListCompleteResult, err error) {
	items := make([]LogAnalyticsQueryPackQuery, 0)

	page, err := c.QueriesList(ctx, id, options)
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

	out := QueriesListCompleteResult{
		Items: items,
	}
	return out, nil
}
