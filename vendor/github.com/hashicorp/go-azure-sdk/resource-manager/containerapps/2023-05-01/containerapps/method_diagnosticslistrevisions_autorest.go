package containerapps

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

type DiagnosticsListRevisionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Revision

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (DiagnosticsListRevisionsOperationResponse, error)
}

type DiagnosticsListRevisionsCompleteResult struct {
	Items []Revision
}

func (r DiagnosticsListRevisionsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r DiagnosticsListRevisionsOperationResponse) LoadMore(ctx context.Context) (resp DiagnosticsListRevisionsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type DiagnosticsListRevisionsOperationOptions struct {
	Filter *string
}

func DefaultDiagnosticsListRevisionsOperationOptions() DiagnosticsListRevisionsOperationOptions {
	return DiagnosticsListRevisionsOperationOptions{}
}

func (o DiagnosticsListRevisionsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o DiagnosticsListRevisionsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// DiagnosticsListRevisions ...
func (c ContainerAppsClient) DiagnosticsListRevisions(ctx context.Context, id ContainerAppId, options DiagnosticsListRevisionsOperationOptions) (resp DiagnosticsListRevisionsOperationResponse, err error) {
	req, err := c.preparerForDiagnosticsListRevisions(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsListRevisions", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsListRevisions", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForDiagnosticsListRevisions(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsListRevisions", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForDiagnosticsListRevisions prepares the DiagnosticsListRevisions request.
func (c ContainerAppsClient) preparerForDiagnosticsListRevisions(ctx context.Context, id ContainerAppId, options DiagnosticsListRevisionsOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/detectorProperties/revisionsApi/revisions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForDiagnosticsListRevisionsWithNextLink prepares the DiagnosticsListRevisions request with the given nextLink token.
func (c ContainerAppsClient) preparerForDiagnosticsListRevisionsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForDiagnosticsListRevisions handles the response to the DiagnosticsListRevisions request. The method always
// closes the http.Response Body.
func (c ContainerAppsClient) responderForDiagnosticsListRevisions(resp *http.Response) (result DiagnosticsListRevisionsOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result DiagnosticsListRevisionsOperationResponse, err error) {
			req, err := c.preparerForDiagnosticsListRevisionsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsListRevisions", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsListRevisions", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForDiagnosticsListRevisions(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsListRevisions", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// DiagnosticsListRevisionsComplete retrieves all of the results into a single object
func (c ContainerAppsClient) DiagnosticsListRevisionsComplete(ctx context.Context, id ContainerAppId, options DiagnosticsListRevisionsOperationOptions) (DiagnosticsListRevisionsCompleteResult, error) {
	return c.DiagnosticsListRevisionsCompleteMatchingPredicate(ctx, id, options, RevisionOperationPredicate{})
}

// DiagnosticsListRevisionsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ContainerAppsClient) DiagnosticsListRevisionsCompleteMatchingPredicate(ctx context.Context, id ContainerAppId, options DiagnosticsListRevisionsOperationOptions, predicate RevisionOperationPredicate) (resp DiagnosticsListRevisionsCompleteResult, err error) {
	items := make([]Revision, 0)

	page, err := c.DiagnosticsListRevisions(ctx, id, options)
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

	out := DiagnosticsListRevisionsCompleteResult{
		Items: items,
	}
	return out, nil
}
