package workspaces

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

type WorkspaceFeaturesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AmlUserFeature

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (WorkspaceFeaturesListOperationResponse, error)
}

type WorkspaceFeaturesListCompleteResult struct {
	Items []AmlUserFeature
}

func (r WorkspaceFeaturesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r WorkspaceFeaturesListOperationResponse) LoadMore(ctx context.Context) (resp WorkspaceFeaturesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// WorkspaceFeaturesList ...
func (c WorkspacesClient) WorkspaceFeaturesList(ctx context.Context, id WorkspaceId) (resp WorkspaceFeaturesListOperationResponse, err error) {
	req, err := c.preparerForWorkspaceFeaturesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "WorkspaceFeaturesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "WorkspaceFeaturesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForWorkspaceFeaturesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "WorkspaceFeaturesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// WorkspaceFeaturesListComplete retrieves all of the results into a single object
func (c WorkspacesClient) WorkspaceFeaturesListComplete(ctx context.Context, id WorkspaceId) (WorkspaceFeaturesListCompleteResult, error) {
	return c.WorkspaceFeaturesListCompleteMatchingPredicate(ctx, id, AmlUserFeatureOperationPredicate{})
}

// WorkspaceFeaturesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c WorkspacesClient) WorkspaceFeaturesListCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, predicate AmlUserFeatureOperationPredicate) (resp WorkspaceFeaturesListCompleteResult, err error) {
	items := make([]AmlUserFeature, 0)

	page, err := c.WorkspaceFeaturesList(ctx, id)
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

	out := WorkspaceFeaturesListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForWorkspaceFeaturesList prepares the WorkspaceFeaturesList request.
func (c WorkspacesClient) preparerForWorkspaceFeaturesList(ctx context.Context, id WorkspaceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/features", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForWorkspaceFeaturesListWithNextLink prepares the WorkspaceFeaturesList request with the given nextLink token.
func (c WorkspacesClient) preparerForWorkspaceFeaturesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForWorkspaceFeaturesList handles the response to the WorkspaceFeaturesList request. The method always
// closes the http.Response Body.
func (c WorkspacesClient) responderForWorkspaceFeaturesList(resp *http.Response) (result WorkspaceFeaturesListOperationResponse, err error) {
	type page struct {
		Values   []AmlUserFeature `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result WorkspaceFeaturesListOperationResponse, err error) {
			req, err := c.preparerForWorkspaceFeaturesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "WorkspaceFeaturesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "WorkspaceFeaturesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForWorkspaceFeaturesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "workspaces.WorkspacesClient", "WorkspaceFeaturesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
