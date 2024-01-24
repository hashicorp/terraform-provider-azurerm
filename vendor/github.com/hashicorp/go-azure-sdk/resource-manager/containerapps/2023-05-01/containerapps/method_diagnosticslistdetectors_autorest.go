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

type DiagnosticsListDetectorsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Diagnostics

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (DiagnosticsListDetectorsOperationResponse, error)
}

type DiagnosticsListDetectorsCompleteResult struct {
	Items []Diagnostics
}

func (r DiagnosticsListDetectorsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r DiagnosticsListDetectorsOperationResponse) LoadMore(ctx context.Context) (resp DiagnosticsListDetectorsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// DiagnosticsListDetectors ...
func (c ContainerAppsClient) DiagnosticsListDetectors(ctx context.Context, id ContainerAppId) (resp DiagnosticsListDetectorsOperationResponse, err error) {
	req, err := c.preparerForDiagnosticsListDetectors(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsListDetectors", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsListDetectors", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForDiagnosticsListDetectors(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsListDetectors", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForDiagnosticsListDetectors prepares the DiagnosticsListDetectors request.
func (c ContainerAppsClient) preparerForDiagnosticsListDetectors(ctx context.Context, id ContainerAppId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/detectors", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForDiagnosticsListDetectorsWithNextLink prepares the DiagnosticsListDetectors request with the given nextLink token.
func (c ContainerAppsClient) preparerForDiagnosticsListDetectorsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForDiagnosticsListDetectors handles the response to the DiagnosticsListDetectors request. The method always
// closes the http.Response Body.
func (c ContainerAppsClient) responderForDiagnosticsListDetectors(resp *http.Response) (result DiagnosticsListDetectorsOperationResponse, err error) {
	type page struct {
		Values   []Diagnostics `json:"value"`
		NextLink *string       `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result DiagnosticsListDetectorsOperationResponse, err error) {
			req, err := c.preparerForDiagnosticsListDetectorsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsListDetectors", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsListDetectors", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForDiagnosticsListDetectors(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "DiagnosticsListDetectors", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// DiagnosticsListDetectorsComplete retrieves all of the results into a single object
func (c ContainerAppsClient) DiagnosticsListDetectorsComplete(ctx context.Context, id ContainerAppId) (DiagnosticsListDetectorsCompleteResult, error) {
	return c.DiagnosticsListDetectorsCompleteMatchingPredicate(ctx, id, DiagnosticsOperationPredicate{})
}

// DiagnosticsListDetectorsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ContainerAppsClient) DiagnosticsListDetectorsCompleteMatchingPredicate(ctx context.Context, id ContainerAppId, predicate DiagnosticsOperationPredicate) (resp DiagnosticsListDetectorsCompleteResult, err error) {
	items := make([]Diagnostics, 0)

	page, err := c.DiagnosticsListDetectors(ctx, id)
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

	out := DiagnosticsListDetectorsCompleteResult{
		Items: items,
	}
	return out, nil
}
