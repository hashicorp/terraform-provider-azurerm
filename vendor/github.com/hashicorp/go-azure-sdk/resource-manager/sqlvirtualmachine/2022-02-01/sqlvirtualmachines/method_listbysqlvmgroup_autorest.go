package sqlvirtualmachines

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

type ListBySqlVmGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SqlVirtualMachine

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListBySqlVmGroupOperationResponse, error)
}

type ListBySqlVmGroupCompleteResult struct {
	Items []SqlVirtualMachine
}

func (r ListBySqlVmGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListBySqlVmGroupOperationResponse) LoadMore(ctx context.Context) (resp ListBySqlVmGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListBySqlVmGroup ...
func (c SqlVirtualMachinesClient) ListBySqlVmGroup(ctx context.Context, id SqlVirtualMachineGroupId) (resp ListBySqlVmGroupOperationResponse, err error) {
	req, err := c.preparerForListBySqlVmGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sqlvirtualmachines.SqlVirtualMachinesClient", "ListBySqlVmGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "sqlvirtualmachines.SqlVirtualMachinesClient", "ListBySqlVmGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListBySqlVmGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sqlvirtualmachines.SqlVirtualMachinesClient", "ListBySqlVmGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListBySqlVmGroupComplete retrieves all of the results into a single object
func (c SqlVirtualMachinesClient) ListBySqlVmGroupComplete(ctx context.Context, id SqlVirtualMachineGroupId) (ListBySqlVmGroupCompleteResult, error) {
	return c.ListBySqlVmGroupCompleteMatchingPredicate(ctx, id, SqlVirtualMachineOperationPredicate{})
}

// ListBySqlVmGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SqlVirtualMachinesClient) ListBySqlVmGroupCompleteMatchingPredicate(ctx context.Context, id SqlVirtualMachineGroupId, predicate SqlVirtualMachineOperationPredicate) (resp ListBySqlVmGroupCompleteResult, err error) {
	items := make([]SqlVirtualMachine, 0)

	page, err := c.ListBySqlVmGroup(ctx, id)
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

	out := ListBySqlVmGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListBySqlVmGroup prepares the ListBySqlVmGroup request.
func (c SqlVirtualMachinesClient) preparerForListBySqlVmGroup(ctx context.Context, id SqlVirtualMachineGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/sqlVirtualMachines", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListBySqlVmGroupWithNextLink prepares the ListBySqlVmGroup request with the given nextLink token.
func (c SqlVirtualMachinesClient) preparerForListBySqlVmGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListBySqlVmGroup handles the response to the ListBySqlVmGroup request. The method always
// closes the http.Response Body.
func (c SqlVirtualMachinesClient) responderForListBySqlVmGroup(resp *http.Response) (result ListBySqlVmGroupOperationResponse, err error) {
	type page struct {
		Values   []SqlVirtualMachine `json:"value"`
		NextLink *string             `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListBySqlVmGroupOperationResponse, err error) {
			req, err := c.preparerForListBySqlVmGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "sqlvirtualmachines.SqlVirtualMachinesClient", "ListBySqlVmGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "sqlvirtualmachines.SqlVirtualMachinesClient", "ListBySqlVmGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListBySqlVmGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "sqlvirtualmachines.SqlVirtualMachinesClient", "ListBySqlVmGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
