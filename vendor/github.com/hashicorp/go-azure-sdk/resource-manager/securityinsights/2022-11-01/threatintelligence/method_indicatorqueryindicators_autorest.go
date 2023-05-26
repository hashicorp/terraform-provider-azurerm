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

type IndicatorQueryIndicatorsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ThreatIntelligenceInformation

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (IndicatorQueryIndicatorsOperationResponse, error)
}

type IndicatorQueryIndicatorsCompleteResult struct {
	Items []ThreatIntelligenceInformation
}

func (r IndicatorQueryIndicatorsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r IndicatorQueryIndicatorsOperationResponse) LoadMore(ctx context.Context) (resp IndicatorQueryIndicatorsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// IndicatorQueryIndicators ...
func (c ThreatIntelligenceClient) IndicatorQueryIndicators(ctx context.Context, id WorkspaceId, input ThreatIntelligenceFilteringCriteria) (resp IndicatorQueryIndicatorsOperationResponse, err error) {
	req, err := c.preparerForIndicatorQueryIndicators(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorQueryIndicators", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorQueryIndicators", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForIndicatorQueryIndicators(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorQueryIndicators", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForIndicatorQueryIndicators prepares the IndicatorQueryIndicators request.
func (c ThreatIntelligenceClient) preparerForIndicatorQueryIndicators(ctx context.Context, id WorkspaceId, input ThreatIntelligenceFilteringCriteria) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.SecurityInsights/threatIntelligence/main/queryIndicators", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForIndicatorQueryIndicatorsWithNextLink prepares the IndicatorQueryIndicators request with the given nextLink token.
func (c ThreatIntelligenceClient) preparerForIndicatorQueryIndicatorsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForIndicatorQueryIndicators handles the response to the IndicatorQueryIndicators request. The method always
// closes the http.Response Body.
func (c ThreatIntelligenceClient) responderForIndicatorQueryIndicators(resp *http.Response) (result IndicatorQueryIndicatorsOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result IndicatorQueryIndicatorsOperationResponse, err error) {
			req, err := c.preparerForIndicatorQueryIndicatorsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorQueryIndicators", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorQueryIndicators", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForIndicatorQueryIndicators(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "threatintelligence.ThreatIntelligenceClient", "IndicatorQueryIndicators", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// IndicatorQueryIndicatorsComplete retrieves all of the results into a single object
func (c ThreatIntelligenceClient) IndicatorQueryIndicatorsComplete(ctx context.Context, id WorkspaceId, input ThreatIntelligenceFilteringCriteria) (IndicatorQueryIndicatorsCompleteResult, error) {
	return c.IndicatorQueryIndicatorsCompleteMatchingPredicate(ctx, id, input, ThreatIntelligenceInformationOperationPredicate{})
}

// IndicatorQueryIndicatorsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ThreatIntelligenceClient) IndicatorQueryIndicatorsCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, input ThreatIntelligenceFilteringCriteria, predicate ThreatIntelligenceInformationOperationPredicate) (resp IndicatorQueryIndicatorsCompleteResult, err error) {
	items := make([]ThreatIntelligenceInformation, 0)

	page, err := c.IndicatorQueryIndicators(ctx, id, input)
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

	out := IndicatorQueryIndicatorsCompleteResult{
		Items: items,
	}
	return out, nil
}
