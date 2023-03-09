package contentkeypolicies

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

type ContentKeyPoliciesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ContentKeyPolicy

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ContentKeyPoliciesListOperationResponse, error)
}

type ContentKeyPoliciesListCompleteResult struct {
	Items []ContentKeyPolicy
}

func (r ContentKeyPoliciesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ContentKeyPoliciesListOperationResponse) LoadMore(ctx context.Context) (resp ContentKeyPoliciesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ContentKeyPoliciesListOperationOptions struct {
	Filter  *string
	Orderby *string
	Top     *int64
}

func DefaultContentKeyPoliciesListOperationOptions() ContentKeyPoliciesListOperationOptions {
	return ContentKeyPoliciesListOperationOptions{}
}

func (o ContentKeyPoliciesListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ContentKeyPoliciesListOperationOptions) toQueryString() map[string]interface{} {
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

// ContentKeyPoliciesList ...
func (c ContentKeyPoliciesClient) ContentKeyPoliciesList(ctx context.Context, id MediaServiceId, options ContentKeyPoliciesListOperationOptions) (resp ContentKeyPoliciesListOperationResponse, err error) {
	req, err := c.preparerForContentKeyPoliciesList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForContentKeyPoliciesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForContentKeyPoliciesList prepares the ContentKeyPoliciesList request.
func (c ContentKeyPoliciesClient) preparerForContentKeyPoliciesList(ctx context.Context, id MediaServiceId, options ContentKeyPoliciesListOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/contentKeyPolicies", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForContentKeyPoliciesListWithNextLink prepares the ContentKeyPoliciesList request with the given nextLink token.
func (c ContentKeyPoliciesClient) preparerForContentKeyPoliciesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForContentKeyPoliciesList handles the response to the ContentKeyPoliciesList request. The method always
// closes the http.Response Body.
func (c ContentKeyPoliciesClient) responderForContentKeyPoliciesList(resp *http.Response) (result ContentKeyPoliciesListOperationResponse, err error) {
	type page struct {
		Values   []ContentKeyPolicy `json:"value"`
		NextLink *string            `json:"@odata.nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ContentKeyPoliciesListOperationResponse, err error) {
			req, err := c.preparerForContentKeyPoliciesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForContentKeyPoliciesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "contentkeypolicies.ContentKeyPoliciesClient", "ContentKeyPoliciesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ContentKeyPoliciesListComplete retrieves all of the results into a single object
func (c ContentKeyPoliciesClient) ContentKeyPoliciesListComplete(ctx context.Context, id MediaServiceId, options ContentKeyPoliciesListOperationOptions) (ContentKeyPoliciesListCompleteResult, error) {
	return c.ContentKeyPoliciesListCompleteMatchingPredicate(ctx, id, options, ContentKeyPolicyOperationPredicate{})
}

// ContentKeyPoliciesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ContentKeyPoliciesClient) ContentKeyPoliciesListCompleteMatchingPredicate(ctx context.Context, id MediaServiceId, options ContentKeyPoliciesListOperationOptions, predicate ContentKeyPolicyOperationPredicate) (resp ContentKeyPoliciesListCompleteResult, err error) {
	items := make([]ContentKeyPolicy, 0)

	page, err := c.ContentKeyPoliciesList(ctx, id, options)
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

	out := ContentKeyPoliciesListCompleteResult{
		Items: items,
	}
	return out, nil
}
