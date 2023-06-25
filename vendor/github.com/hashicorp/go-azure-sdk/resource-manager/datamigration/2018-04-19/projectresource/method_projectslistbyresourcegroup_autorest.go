package projectresource

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

type ProjectsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Project

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ProjectsListByResourceGroupOperationResponse, error)
}

type ProjectsListByResourceGroupCompleteResult struct {
	Items []Project
}

func (r ProjectsListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ProjectsListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp ProjectsListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ProjectsListByResourceGroup ...
func (c ProjectResourceClient) ProjectsListByResourceGroup(ctx context.Context, id ServiceId) (resp ProjectsListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForProjectsListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForProjectsListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForProjectsListByResourceGroup prepares the ProjectsListByResourceGroup request.
func (c ProjectResourceClient) preparerForProjectsListByResourceGroup(ctx context.Context, id ServiceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/projects", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForProjectsListByResourceGroupWithNextLink prepares the ProjectsListByResourceGroup request with the given nextLink token.
func (c ProjectResourceClient) preparerForProjectsListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForProjectsListByResourceGroup handles the response to the ProjectsListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ProjectResourceClient) responderForProjectsListByResourceGroup(resp *http.Response) (result ProjectsListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []Project `json:"value"`
		NextLink *string   `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ProjectsListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForProjectsListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForProjectsListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ProjectsListByResourceGroupComplete retrieves all of the results into a single object
func (c ProjectResourceClient) ProjectsListByResourceGroupComplete(ctx context.Context, id ServiceId) (ProjectsListByResourceGroupCompleteResult, error) {
	return c.ProjectsListByResourceGroupCompleteMatchingPredicate(ctx, id, ProjectOperationPredicate{})
}

// ProjectsListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ProjectResourceClient) ProjectsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id ServiceId, predicate ProjectOperationPredicate) (resp ProjectsListByResourceGroupCompleteResult, err error) {
	items := make([]Project, 0)

	page, err := c.ProjectsListByResourceGroup(ctx, id)
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

	out := ProjectsListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
