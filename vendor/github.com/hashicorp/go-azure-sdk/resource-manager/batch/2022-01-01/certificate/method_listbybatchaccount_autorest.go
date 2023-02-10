package certificate

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

type ListByBatchAccountOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Certificate

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByBatchAccountOperationResponse, error)
}

type ListByBatchAccountCompleteResult struct {
	Items []Certificate
}

func (r ListByBatchAccountOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByBatchAccountOperationResponse) LoadMore(ctx context.Context) (resp ListByBatchAccountOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByBatchAccountOperationOptions struct {
	Filter     *string
	Maxresults *int64
	Select     *string
}

func DefaultListByBatchAccountOperationOptions() ListByBatchAccountOperationOptions {
	return ListByBatchAccountOperationOptions{}
}

func (o ListByBatchAccountOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListByBatchAccountOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Maxresults != nil {
		out["maxresults"] = *o.Maxresults
	}

	if o.Select != nil {
		out["$select"] = *o.Select
	}

	return out
}

// ListByBatchAccount ...
func (c CertificateClient) ListByBatchAccount(ctx context.Context, id BatchAccountId, options ListByBatchAccountOperationOptions) (resp ListByBatchAccountOperationResponse, err error) {
	req, err := c.preparerForListByBatchAccount(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificate.CertificateClient", "ListByBatchAccount", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificate.CertificateClient", "ListByBatchAccount", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByBatchAccount(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "certificate.CertificateClient", "ListByBatchAccount", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByBatchAccount prepares the ListByBatchAccount request.
func (c CertificateClient) preparerForListByBatchAccount(ctx context.Context, id BatchAccountId, options ListByBatchAccountOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/certificates", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByBatchAccountWithNextLink prepares the ListByBatchAccount request with the given nextLink token.
func (c CertificateClient) preparerForListByBatchAccountWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByBatchAccount handles the response to the ListByBatchAccount request. The method always
// closes the http.Response Body.
func (c CertificateClient) responderForListByBatchAccount(resp *http.Response) (result ListByBatchAccountOperationResponse, err error) {
	type page struct {
		Values   []Certificate `json:"value"`
		NextLink *string       `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByBatchAccountOperationResponse, err error) {
			req, err := c.preparerForListByBatchAccountWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "certificate.CertificateClient", "ListByBatchAccount", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "certificate.CertificateClient", "ListByBatchAccount", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByBatchAccount(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "certificate.CertificateClient", "ListByBatchAccount", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByBatchAccountComplete retrieves all of the results into a single object
func (c CertificateClient) ListByBatchAccountComplete(ctx context.Context, id BatchAccountId, options ListByBatchAccountOperationOptions) (ListByBatchAccountCompleteResult, error) {
	return c.ListByBatchAccountCompleteMatchingPredicate(ctx, id, options, CertificateOperationPredicate{})
}

// ListByBatchAccountCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c CertificateClient) ListByBatchAccountCompleteMatchingPredicate(ctx context.Context, id BatchAccountId, options ListByBatchAccountOperationOptions, predicate CertificateOperationPredicate) (resp ListByBatchAccountCompleteResult, err error) {
	items := make([]Certificate, 0)

	page, err := c.ListByBatchAccount(ctx, id, options)
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

	out := ListByBatchAccountCompleteResult{
		Items: items,
	}
	return out, nil
}
