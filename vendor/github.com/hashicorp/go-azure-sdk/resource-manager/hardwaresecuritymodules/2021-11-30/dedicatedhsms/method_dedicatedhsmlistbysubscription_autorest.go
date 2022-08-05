package dedicatedhsms

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

type DedicatedHsmListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DedicatedHsm

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (DedicatedHsmListBySubscriptionOperationResponse, error)
}

type DedicatedHsmListBySubscriptionCompleteResult struct {
	Items []DedicatedHsm
}

func (r DedicatedHsmListBySubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r DedicatedHsmListBySubscriptionOperationResponse) LoadMore(ctx context.Context) (resp DedicatedHsmListBySubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type DedicatedHsmListBySubscriptionOperationOptions struct {
	Top *int64
}

func DefaultDedicatedHsmListBySubscriptionOperationOptions() DedicatedHsmListBySubscriptionOperationOptions {
	return DedicatedHsmListBySubscriptionOperationOptions{}
}

func (o DedicatedHsmListBySubscriptionOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o DedicatedHsmListBySubscriptionOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// DedicatedHsmListBySubscription ...
func (c DedicatedHsmsClient) DedicatedHsmListBySubscription(ctx context.Context, id commonids.SubscriptionId, options DedicatedHsmListBySubscriptionOperationOptions) (resp DedicatedHsmListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForDedicatedHsmListBySubscription(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForDedicatedHsmListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// DedicatedHsmListBySubscriptionComplete retrieves all of the results into a single object
func (c DedicatedHsmsClient) DedicatedHsmListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId, options DedicatedHsmListBySubscriptionOperationOptions) (DedicatedHsmListBySubscriptionCompleteResult, error) {
	return c.DedicatedHsmListBySubscriptionCompleteMatchingPredicate(ctx, id, options, DedicatedHsmOperationPredicate{})
}

// DedicatedHsmListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DedicatedHsmsClient) DedicatedHsmListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options DedicatedHsmListBySubscriptionOperationOptions, predicate DedicatedHsmOperationPredicate) (resp DedicatedHsmListBySubscriptionCompleteResult, err error) {
	items := make([]DedicatedHsm, 0)

	page, err := c.DedicatedHsmListBySubscription(ctx, id, options)
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

	out := DedicatedHsmListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForDedicatedHsmListBySubscription prepares the DedicatedHsmListBySubscription request.
func (c DedicatedHsmsClient) preparerForDedicatedHsmListBySubscription(ctx context.Context, id commonids.SubscriptionId, options DedicatedHsmListBySubscriptionOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.HardwareSecurityModules/dedicatedHSMs", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForDedicatedHsmListBySubscriptionWithNextLink prepares the DedicatedHsmListBySubscription request with the given nextLink token.
func (c DedicatedHsmsClient) preparerForDedicatedHsmListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForDedicatedHsmListBySubscription handles the response to the DedicatedHsmListBySubscription request. The method always
// closes the http.Response Body.
func (c DedicatedHsmsClient) responderForDedicatedHsmListBySubscription(resp *http.Response) (result DedicatedHsmListBySubscriptionOperationResponse, err error) {
	type page struct {
		Values   []DedicatedHsm `json:"value"`
		NextLink *string        `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result DedicatedHsmListBySubscriptionOperationResponse, err error) {
			req, err := c.preparerForDedicatedHsmListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForDedicatedHsmListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
