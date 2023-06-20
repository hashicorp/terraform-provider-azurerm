package monitors

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

type ListUserRolesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]UserRoleResponse

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListUserRolesOperationResponse, error)
}

type ListUserRolesCompleteResult struct {
	Items []UserRoleResponse
}

func (r ListUserRolesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListUserRolesOperationResponse) LoadMore(ctx context.Context) (resp ListUserRolesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListUserRoles ...
func (c MonitorsClient) ListUserRoles(ctx context.Context, id MonitorId, input UserRoleRequest) (resp ListUserRolesOperationResponse, err error) {
	req, err := c.preparerForListUserRoles(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListUserRoles", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListUserRoles", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListUserRoles(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListUserRoles", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListUserRoles prepares the ListUserRoles request.
func (c MonitorsClient) preparerForListUserRoles(ctx context.Context, id MonitorId, input UserRoleRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listUserRoles", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListUserRolesWithNextLink prepares the ListUserRoles request with the given nextLink token.
func (c MonitorsClient) preparerForListUserRolesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListUserRoles handles the response to the ListUserRoles request. The method always
// closes the http.Response Body.
func (c MonitorsClient) responderForListUserRoles(resp *http.Response) (result ListUserRolesOperationResponse, err error) {
	type page struct {
		Values   []UserRoleResponse `json:"value"`
		NextLink *string            `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListUserRolesOperationResponse, err error) {
			req, err := c.preparerForListUserRolesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListUserRoles", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListUserRoles", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListUserRoles(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "ListUserRoles", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListUserRolesComplete retrieves all of the results into a single object
func (c MonitorsClient) ListUserRolesComplete(ctx context.Context, id MonitorId, input UserRoleRequest) (ListUserRolesCompleteResult, error) {
	return c.ListUserRolesCompleteMatchingPredicate(ctx, id, input, UserRoleResponseOperationPredicate{})
}

// ListUserRolesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c MonitorsClient) ListUserRolesCompleteMatchingPredicate(ctx context.Context, id MonitorId, input UserRoleRequest, predicate UserRoleResponseOperationPredicate) (resp ListUserRolesCompleteResult, err error) {
	items := make([]UserRoleResponse, 0)

	page, err := c.ListUserRoles(ctx, id, input)
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

	out := ListUserRolesCompleteResult{
		Items: items,
	}
	return out, nil
}
