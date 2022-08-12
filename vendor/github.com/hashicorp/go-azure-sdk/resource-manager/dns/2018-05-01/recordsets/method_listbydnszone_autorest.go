package recordsets

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

type ListByDnsZoneOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]RecordSet

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByDnsZoneOperationResponse, error)
}

type ListByDnsZoneCompleteResult struct {
	Items []RecordSet
}

func (r ListByDnsZoneOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByDnsZoneOperationResponse) LoadMore(ctx context.Context) (resp ListByDnsZoneOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByDnsZoneOperationOptions struct {
	Recordsetnamesuffix *string
	Top                 *int64
}

func DefaultListByDnsZoneOperationOptions() ListByDnsZoneOperationOptions {
	return ListByDnsZoneOperationOptions{}
}

func (o ListByDnsZoneOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListByDnsZoneOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Recordsetnamesuffix != nil {
		out["$recordsetnamesuffix"] = *o.Recordsetnamesuffix
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListByDnsZone ...
func (c RecordSetsClient) ListByDnsZone(ctx context.Context, id DnsZoneId, options ListByDnsZoneOperationOptions) (resp ListByDnsZoneOperationResponse, err error) {
	req, err := c.preparerForListByDnsZone(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "recordsets.RecordSetsClient", "ListByDnsZone", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "recordsets.RecordSetsClient", "ListByDnsZone", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByDnsZone(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "recordsets.RecordSetsClient", "ListByDnsZone", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByDnsZoneComplete retrieves all of the results into a single object
func (c RecordSetsClient) ListByDnsZoneComplete(ctx context.Context, id DnsZoneId, options ListByDnsZoneOperationOptions) (ListByDnsZoneCompleteResult, error) {
	return c.ListByDnsZoneCompleteMatchingPredicate(ctx, id, options, RecordSetOperationPredicate{})
}

// ListByDnsZoneCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c RecordSetsClient) ListByDnsZoneCompleteMatchingPredicate(ctx context.Context, id DnsZoneId, options ListByDnsZoneOperationOptions, predicate RecordSetOperationPredicate) (resp ListByDnsZoneCompleteResult, err error) {
	items := make([]RecordSet, 0)

	page, err := c.ListByDnsZone(ctx, id, options)
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

	out := ListByDnsZoneCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByDnsZone prepares the ListByDnsZone request.
func (c RecordSetsClient) preparerForListByDnsZone(ctx context.Context, id DnsZoneId, options ListByDnsZoneOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/recordsets", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByDnsZoneWithNextLink prepares the ListByDnsZone request with the given nextLink token.
func (c RecordSetsClient) preparerForListByDnsZoneWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByDnsZone handles the response to the ListByDnsZone request. The method always
// closes the http.Response Body.
func (c RecordSetsClient) responderForListByDnsZone(resp *http.Response) (result ListByDnsZoneOperationResponse, err error) {
	type page struct {
		Values   []RecordSet `json:"value"`
		NextLink *string     `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByDnsZoneOperationResponse, err error) {
			req, err := c.preparerForListByDnsZoneWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "recordsets.RecordSetsClient", "ListByDnsZone", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "recordsets.RecordSetsClient", "ListByDnsZone", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByDnsZone(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "recordsets.RecordSetsClient", "ListByDnsZone", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
