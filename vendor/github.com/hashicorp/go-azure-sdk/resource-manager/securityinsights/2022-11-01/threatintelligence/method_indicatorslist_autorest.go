package threatintelligence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IndicatorsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ThreatIntelligenceInformation

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (IndicatorsListOperationResponse, error)
}

type IndicatorsListCompleteResult struct {
	Items []ThreatIntelligenceInformation
}

func (r IndicatorsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r IndicatorsListOperationResponse) LoadMore(ctx context.Context) (resp IndicatorsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type IndicatorsListOperationOptions struct {
	Filter  *string
	Orderby *string
	Top     *int64
}

func DefaultIndicatorsListOperationOptions() IndicatorsListOperationOptions {
	return IndicatorsListOperationOptions{}
}

func (o IndicatorsListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o IndicatorsListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Orderby != nil {
		out["$orderby"] = *o.Orderby
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// IndicatorsList ...
func (c ThreatIntelligenceClient) IndicatorsList(ctx context.Context, id WorkspaceId, options IndicatorsListOperationOptions) (resp IndicatorsListOperationResponse, err error) {
	req, err := c.preparerForIndicatorsList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForIndicatorsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForIndicatorsList prepares the IndicatorsList request.
func (c ThreatIntelligenceClient) preparerForIndicatorsList(ctx context.Context, id WorkspaceId, options IndicatorsListOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.SecurityInsights/threatIntelligence/main/indicators", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForIndicatorsListWithNextLink prepares the IndicatorsList request with the given nextLink token.
func (c ThreatIntelligenceClient) preparerForIndicatorsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForIndicatorsList handles the response to the IndicatorsList request. The method always
// closes the http.Response Body.
func (c ThreatIntelligenceClient) responderForIndicatorsList(resp *http.Response) (result IndicatorsListOperationResponse, err error) {
	type page struct {
		Values   []json.RawMessage `json:"value"`
		NextLink *string           `json:"nextLink"`
	}
	var respObj page
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&respObj),
		autorest.ByClosing())
	result.HttpResponse = resp
	temp := make([]ThreatIntelligenceInformation, 0)
	for i, v := range respObj.Values {
		val, err := unmarshalThreatIntelligenceInformationImplementation(v)
		if err != nil {
			err = fmt.Errorf("unmarshalling item %d for ThreatIntelligenceInformation (%q): %+v", i, v, err)
			return result, err
		}
		temp = append(temp, val)
	}
	result.Model = &temp
	result.nextLink = respObj.NextLink
	if respObj.NextLink != nil {
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result IndicatorsListOperationResponse, err error) {
			req, err := c.preparerForIndicatorsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForIndicatorsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// IndicatorsListComplete retrieves all of the results into a single object
func (c ThreatIntelligenceClient) IndicatorsListComplete(ctx context.Context, id WorkspaceId, options IndicatorsListOperationOptions) (IndicatorsListCompleteResult, error) {
	return c.IndicatorsListCompleteMatchingPredicate(ctx, id, options, ThreatIntelligenceInformationOperationPredicate{})
}

// IndicatorsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ThreatIntelligenceClient) IndicatorsListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, options IndicatorsListOperationOptions, predicate ThreatIntelligenceInformationOperationPredicate) (resp IndicatorsListCompleteResult, err error) {
	items := make([]ThreatIntelligenceInformation, 0)

	page, err := c.IndicatorsList(ctx, id, options)
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

	out := IndicatorsListCompleteResult{
		Items: items,
	}
	return out, nil
}
