package containerappsrevisions

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

type ListRevisionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Revision

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListRevisionsOperationResponse, error)
}

type ListRevisionsCompleteResult struct {
	Items []Revision
}

func (r ListRevisionsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListRevisionsOperationResponse) LoadMore(ctx context.Context) (resp ListRevisionsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListRevisionsOperationOptions struct {
	Filter *string
}

func DefaultListRevisionsOperationOptions() ListRevisionsOperationOptions {
	return ListRevisionsOperationOptions{}
}

func (o ListRevisionsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListRevisionsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// ListRevisions ...
func (c ContainerAppsRevisionsClient) ListRevisions(ctx context.Context, id ContainerAppId, options ListRevisionsOperationOptions) (resp ListRevisionsOperationResponse, err error) {
	req, err := c.preparerForListRevisions(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "ListRevisions", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "ListRevisions", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListRevisions(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "ListRevisions", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListRevisions prepares the ListRevisions request.
func (c ContainerAppsRevisionsClient) preparerForListRevisions(ctx context.Context, id ContainerAppId, options ListRevisionsOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/revisions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListRevisionsWithNextLink prepares the ListRevisions request with the given nextLink token.
func (c ContainerAppsRevisionsClient) preparerForListRevisionsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListRevisions handles the response to the ListRevisions request. The method always
// closes the http.Response Body.
func (c ContainerAppsRevisionsClient) responderForListRevisions(resp *http.Response) (result ListRevisionsOperationResponse, err error) {
	type page struct {
		Values   []Revision `json:"value"`
		NextLink *string    `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListRevisionsOperationResponse, err error) {
			req, err := c.preparerForListRevisionsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "ListRevisions", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "ListRevisions", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListRevisions(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "ListRevisions", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListRevisionsComplete retrieves all of the results into a single object
func (c ContainerAppsRevisionsClient) ListRevisionsComplete(ctx context.Context, id ContainerAppId, options ListRevisionsOperationOptions) (ListRevisionsCompleteResult, error) {
	return c.ListRevisionsCompleteMatchingPredicate(ctx, id, options, RevisionOperationPredicate{})
}

// ListRevisionsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ContainerAppsRevisionsClient) ListRevisionsCompleteMatchingPredicate(ctx context.Context, id ContainerAppId, options ListRevisionsOperationOptions, predicate RevisionOperationPredicate) (resp ListRevisionsCompleteResult, err error) {
	items := make([]Revision, 0)

	page, err := c.ListRevisions(ctx, id, options)
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

	out := ListRevisionsCompleteResult{
		Items: items,
	}
	return out, nil
}
