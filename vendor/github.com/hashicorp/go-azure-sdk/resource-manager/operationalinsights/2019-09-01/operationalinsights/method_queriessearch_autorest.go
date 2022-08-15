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

type QueriesSearchOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]LogAnalyticsQueryPackQuery

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (QueriesSearchOperationResponse, error)
}

type QueriesSearchCompleteResult struct {
	Items []LogAnalyticsQueryPackQuery
}

func (r QueriesSearchOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r QueriesSearchOperationResponse) LoadMore(ctx context.Context) (resp QueriesSearchOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type QueriesSearchOperationOptions struct {
	IncludeBody *bool
	Top         *int64
}

func DefaultQueriesSearchOperationOptions() QueriesSearchOperationOptions {
	return QueriesSearchOperationOptions{}
}

func (o QueriesSearchOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o QueriesSearchOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IncludeBody != nil {
		out["includeBody"] = *o.IncludeBody
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// QueriesSearch ...
func (c OperationalInsightsClient) QueriesSearch(ctx context.Context, id QueryPackId, input LogAnalyticsQueryPackQuerySearchProperties, options QueriesSearchOperationOptions) (resp QueriesSearchOperationResponse, err error) {
	req, err := c.preparerForQueriesSearch(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesSearch", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesSearch", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForQueriesSearch(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesSearch", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForQueriesSearch prepares the QueriesSearch request.
func (c OperationalInsightsClient) preparerForQueriesSearch(ctx context.Context, id QueryPackId, input LogAnalyticsQueryPackQuerySearchProperties, options QueriesSearchOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/queries/search", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForQueriesSearchWithNextLink prepares the QueriesSearch request with the given nextLink token.
func (c OperationalInsightsClient) preparerForQueriesSearchWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForQueriesSearch handles the response to the QueriesSearch request. The method always
// closes the http.Response Body.
func (c OperationalInsightsClient) responderForQueriesSearch(resp *http.Response) (result QueriesSearchOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result QueriesSearchOperationResponse, err error) {
			req, err := c.preparerForQueriesSearchWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesSearch", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesSearch", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForQueriesSearch(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationalinsights.OperationalInsightsClient", "QueriesSearch", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// QueriesSearchComplete retrieves all of the results into a single object
func (c OperationalInsightsClient) QueriesSearchComplete(ctx context.Context, id QueryPackId, input LogAnalyticsQueryPackQuerySearchProperties, options QueriesSearchOperationOptions) (QueriesSearchCompleteResult, error) {
	return c.QueriesSearchCompleteMatchingPredicate(ctx, id, input, options, LogAnalyticsQueryPackQueryOperationPredicate{})
}

// QueriesSearchCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c OperationalInsightsClient) QueriesSearchCompleteMatchingPredicate(ctx context.Context, id QueryPackId, input LogAnalyticsQueryPackQuerySearchProperties, options QueriesSearchOperationOptions, predicate LogAnalyticsQueryPackQueryOperationPredicate) (resp QueriesSearchCompleteResult, err error) {
	items := make([]LogAnalyticsQueryPackQuery, 0)

	page, err := c.QueriesSearch(ctx, id, input, options)
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

	out := QueriesSearchCompleteResult{
		Items: items,
	}
	return out, nil
}
