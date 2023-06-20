package schema

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

type GlobalSchemaListByServiceOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]GlobalSchemaContract

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (GlobalSchemaListByServiceOperationResponse, error)
}

type GlobalSchemaListByServiceCompleteResult struct {
	Items []GlobalSchemaContract
}

func (r GlobalSchemaListByServiceOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r GlobalSchemaListByServiceOperationResponse) LoadMore(ctx context.Context) (resp GlobalSchemaListByServiceOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type GlobalSchemaListByServiceOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultGlobalSchemaListByServiceOperationOptions() GlobalSchemaListByServiceOperationOptions {
	return GlobalSchemaListByServiceOperationOptions{}
}

func (o GlobalSchemaListByServiceOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o GlobalSchemaListByServiceOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Skip != nil {
		out["$skip"] = *o.Skip
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// GlobalSchemaListByService ...
func (c SchemaClient) GlobalSchemaListByService(ctx context.Context, id ServiceId, options GlobalSchemaListByServiceOperationOptions) (resp GlobalSchemaListByServiceOperationResponse, err error) {
	req, err := c.preparerForGlobalSchemaListByService(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaListByService", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaListByService", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForGlobalSchemaListByService(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaListByService", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForGlobalSchemaListByService prepares the GlobalSchemaListByService request.
func (c SchemaClient) preparerForGlobalSchemaListByService(ctx context.Context, id ServiceId, options GlobalSchemaListByServiceOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/schemas", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForGlobalSchemaListByServiceWithNextLink prepares the GlobalSchemaListByService request with the given nextLink token.
func (c SchemaClient) preparerForGlobalSchemaListByServiceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForGlobalSchemaListByService handles the response to the GlobalSchemaListByService request. The method always
// closes the http.Response Body.
func (c SchemaClient) responderForGlobalSchemaListByService(resp *http.Response) (result GlobalSchemaListByServiceOperationResponse, err error) {
	type page struct {
		Values   []GlobalSchemaContract `json:"value"`
		NextLink *string                `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result GlobalSchemaListByServiceOperationResponse, err error) {
			req, err := c.preparerForGlobalSchemaListByServiceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaListByService", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaListByService", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForGlobalSchemaListByService(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaListByService", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// GlobalSchemaListByServiceComplete retrieves all of the results into a single object
func (c SchemaClient) GlobalSchemaListByServiceComplete(ctx context.Context, id ServiceId, options GlobalSchemaListByServiceOperationOptions) (GlobalSchemaListByServiceCompleteResult, error) {
	return c.GlobalSchemaListByServiceCompleteMatchingPredicate(ctx, id, options, GlobalSchemaContractOperationPredicate{})
}

// GlobalSchemaListByServiceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SchemaClient) GlobalSchemaListByServiceCompleteMatchingPredicate(ctx context.Context, id ServiceId, options GlobalSchemaListByServiceOperationOptions, predicate GlobalSchemaContractOperationPredicate) (resp GlobalSchemaListByServiceCompleteResult, err error) {
	items := make([]GlobalSchemaContract, 0)

	page, err := c.GlobalSchemaListByService(ctx, id, options)
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

	out := GlobalSchemaListByServiceCompleteResult{
		Items: items,
	}
	return out, nil
}
