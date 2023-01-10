package virtualmachines

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

type ListByLocationOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]VirtualMachine

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByLocationOperationResponse, error)
}

type ListByLocationCompleteResult struct {
	Items []VirtualMachine
}

func (r ListByLocationOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByLocationOperationResponse) LoadMore(ctx context.Context) (resp ListByLocationOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByLocation ...
func (c VirtualMachinesClient) ListByLocation(ctx context.Context, id LocationId) (resp ListByLocationOperationResponse, err error) {
	req, err := c.preparerForListByLocation(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "ListByLocation", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "ListByLocation", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByLocation(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "ListByLocation", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByLocation prepares the ListByLocation request.
func (c VirtualMachinesClient) preparerForListByLocation(ctx context.Context, id LocationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/virtualMachines", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByLocationWithNextLink prepares the ListByLocation request with the given nextLink token.
func (c VirtualMachinesClient) preparerForListByLocationWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByLocation handles the response to the ListByLocation request. The method always
// closes the http.Response Body.
func (c VirtualMachinesClient) responderForListByLocation(resp *http.Response) (result ListByLocationOperationResponse, err error) {
	type page struct {
		Values   []VirtualMachine `json:"value"`
		NextLink *string          `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByLocationOperationResponse, err error) {
			req, err := c.preparerForListByLocationWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "ListByLocation", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "ListByLocation", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByLocation(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "ListByLocation", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByLocationComplete retrieves all of the results into a single object
func (c VirtualMachinesClient) ListByLocationComplete(ctx context.Context, id LocationId) (ListByLocationCompleteResult, error) {
	return c.ListByLocationCompleteMatchingPredicate(ctx, id, VirtualMachineOperationPredicate{})
}

// ListByLocationCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c VirtualMachinesClient) ListByLocationCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate VirtualMachineOperationPredicate) (resp ListByLocationCompleteResult, err error) {
	items := make([]VirtualMachine, 0)

	page, err := c.ListByLocation(ctx, id)
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

	out := ListByLocationCompleteResult{
		Items: items,
	}
	return out, nil
}
