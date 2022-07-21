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

type DedicatedHsmListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DedicatedHsm

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (DedicatedHsmListByResourceGroupOperationResponse, error)
}

type DedicatedHsmListByResourceGroupCompleteResult struct {
	Items []DedicatedHsm
}

func (r DedicatedHsmListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r DedicatedHsmListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp DedicatedHsmListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type DedicatedHsmListByResourceGroupOperationOptions struct {
	Top *int64
}

func DefaultDedicatedHsmListByResourceGroupOperationOptions() DedicatedHsmListByResourceGroupOperationOptions {
	return DedicatedHsmListByResourceGroupOperationOptions{}
}

func (o DedicatedHsmListByResourceGroupOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o DedicatedHsmListByResourceGroupOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// DedicatedHsmListByResourceGroup ...
func (c DedicatedHsmsClient) DedicatedHsmListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options DedicatedHsmListByResourceGroupOperationOptions) (resp DedicatedHsmListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForDedicatedHsmListByResourceGroup(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForDedicatedHsmListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// DedicatedHsmListByResourceGroupComplete retrieves all of the results into a single object
func (c DedicatedHsmsClient) DedicatedHsmListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId, options DedicatedHsmListByResourceGroupOperationOptions) (DedicatedHsmListByResourceGroupCompleteResult, error) {
	return c.DedicatedHsmListByResourceGroupCompleteMatchingPredicate(ctx, id, options, DedicatedHsmOperationPredicate{})
}

// DedicatedHsmListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DedicatedHsmsClient) DedicatedHsmListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, options DedicatedHsmListByResourceGroupOperationOptions, predicate DedicatedHsmOperationPredicate) (resp DedicatedHsmListByResourceGroupCompleteResult, err error) {
	items := make([]DedicatedHsm, 0)

	page, err := c.DedicatedHsmListByResourceGroup(ctx, id, options)
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

	out := DedicatedHsmListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForDedicatedHsmListByResourceGroup prepares the DedicatedHsmListByResourceGroup request.
func (c DedicatedHsmsClient) preparerForDedicatedHsmListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId, options DedicatedHsmListByResourceGroupOperationOptions) (*http.Request, error) {
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

// preparerForDedicatedHsmListByResourceGroupWithNextLink prepares the DedicatedHsmListByResourceGroup request with the given nextLink token.
func (c DedicatedHsmsClient) preparerForDedicatedHsmListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForDedicatedHsmListByResourceGroup handles the response to the DedicatedHsmListByResourceGroup request. The method always
// closes the http.Response Body.
func (c DedicatedHsmsClient) responderForDedicatedHsmListByResourceGroup(resp *http.Response) (result DedicatedHsmListByResourceGroupOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result DedicatedHsmListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForDedicatedHsmListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForDedicatedHsmListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "dedicatedhsms.DedicatedHsmsClient", "DedicatedHsmListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
