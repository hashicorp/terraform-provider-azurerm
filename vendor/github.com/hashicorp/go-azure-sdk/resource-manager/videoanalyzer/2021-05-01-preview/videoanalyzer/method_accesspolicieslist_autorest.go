package videoanalyzer

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

type AccessPoliciesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AccessPolicyEntity

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AccessPoliciesListOperationResponse, error)
}

type AccessPoliciesListCompleteResult struct {
	Items []AccessPolicyEntity
}

func (r AccessPoliciesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AccessPoliciesListOperationResponse) LoadMore(ctx context.Context) (resp AccessPoliciesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type AccessPoliciesListOperationOptions struct {
	Top *int64
}

func DefaultAccessPoliciesListOperationOptions() AccessPoliciesListOperationOptions {
	return AccessPoliciesListOperationOptions{}
}

func (o AccessPoliciesListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o AccessPoliciesListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// AccessPoliciesList ...
func (c VideoAnalyzerClient) AccessPoliciesList(ctx context.Context, id VideoAnalyzerId, options AccessPoliciesListOperationOptions) (resp AccessPoliciesListOperationResponse, err error) {
	req, err := c.preparerForAccessPoliciesList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAccessPoliciesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// AccessPoliciesListComplete retrieves all of the results into a single object
func (c VideoAnalyzerClient) AccessPoliciesListComplete(ctx context.Context, id VideoAnalyzerId, options AccessPoliciesListOperationOptions) (AccessPoliciesListCompleteResult, error) {
	return c.AccessPoliciesListCompleteMatchingPredicate(ctx, id, options, AccessPolicyEntityOperationPredicate{})
}

// AccessPoliciesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c VideoAnalyzerClient) AccessPoliciesListCompleteMatchingPredicate(ctx context.Context, id VideoAnalyzerId, options AccessPoliciesListOperationOptions, predicate AccessPolicyEntityOperationPredicate) (resp AccessPoliciesListCompleteResult, err error) {
	items := make([]AccessPolicyEntity, 0)

	page, err := c.AccessPoliciesList(ctx, id, options)
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

	out := AccessPoliciesListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForAccessPoliciesList prepares the AccessPoliciesList request.
func (c VideoAnalyzerClient) preparerForAccessPoliciesList(ctx context.Context, id VideoAnalyzerId, options AccessPoliciesListOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/accessPolicies", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForAccessPoliciesListWithNextLink prepares the AccessPoliciesList request with the given nextLink token.
func (c VideoAnalyzerClient) preparerForAccessPoliciesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAccessPoliciesList handles the response to the AccessPoliciesList request. The method always
// closes the http.Response Body.
func (c VideoAnalyzerClient) responderForAccessPoliciesList(resp *http.Response) (result AccessPoliciesListOperationResponse, err error) {
	type page struct {
		Values   []AccessPolicyEntity `json:"value"`
		NextLink *string              `json:"@nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AccessPoliciesListOperationResponse, err error) {
			req, err := c.preparerForAccessPoliciesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAccessPoliciesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "videoanalyzer.VideoAnalyzerClient", "AccessPoliciesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
