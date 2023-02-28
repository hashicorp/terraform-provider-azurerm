package frontdoors

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

type RulesEnginesListByFrontDoorOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]RulesEngine

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (RulesEnginesListByFrontDoorOperationResponse, error)
}

type RulesEnginesListByFrontDoorCompleteResult struct {
	Items []RulesEngine
}

func (r RulesEnginesListByFrontDoorOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r RulesEnginesListByFrontDoorOperationResponse) LoadMore(ctx context.Context) (resp RulesEnginesListByFrontDoorOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// RulesEnginesListByFrontDoor ...
func (c FrontDoorsClient) RulesEnginesListByFrontDoor(ctx context.Context, id FrontDoorId) (resp RulesEnginesListByFrontDoorOperationResponse, err error) {
	req, err := c.preparerForRulesEnginesListByFrontDoor(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "RulesEnginesListByFrontDoor", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "RulesEnginesListByFrontDoor", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForRulesEnginesListByFrontDoor(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "RulesEnginesListByFrontDoor", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForRulesEnginesListByFrontDoor prepares the RulesEnginesListByFrontDoor request.
func (c FrontDoorsClient) preparerForRulesEnginesListByFrontDoor(ctx context.Context, id FrontDoorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/rulesEngines", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForRulesEnginesListByFrontDoorWithNextLink prepares the RulesEnginesListByFrontDoor request with the given nextLink token.
func (c FrontDoorsClient) preparerForRulesEnginesListByFrontDoorWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForRulesEnginesListByFrontDoor handles the response to the RulesEnginesListByFrontDoor request. The method always
// closes the http.Response Body.
func (c FrontDoorsClient) responderForRulesEnginesListByFrontDoor(resp *http.Response) (result RulesEnginesListByFrontDoorOperationResponse, err error) {
	type page struct {
		Values   []RulesEngine `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result RulesEnginesListByFrontDoorOperationResponse, err error) {
			req, err := c.preparerForRulesEnginesListByFrontDoorWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "RulesEnginesListByFrontDoor", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "RulesEnginesListByFrontDoor", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForRulesEnginesListByFrontDoor(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "frontdoors.FrontDoorsClient", "RulesEnginesListByFrontDoor", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// RulesEnginesListByFrontDoorComplete retrieves all of the results into a single object
func (c FrontDoorsClient) RulesEnginesListByFrontDoorComplete(ctx context.Context, id FrontDoorId) (RulesEnginesListByFrontDoorCompleteResult, error) {
	return c.RulesEnginesListByFrontDoorCompleteMatchingPredicate(ctx, id, RulesEngineOperationPredicate{})
}

// RulesEnginesListByFrontDoorCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c FrontDoorsClient) RulesEnginesListByFrontDoorCompleteMatchingPredicate(ctx context.Context, id FrontDoorId, predicate RulesEngineOperationPredicate) (resp RulesEnginesListByFrontDoorCompleteResult, err error) {
	items := make([]RulesEngine, 0)

	page, err := c.RulesEnginesListByFrontDoor(ctx, id)
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

	out := RulesEnginesListByFrontDoorCompleteResult{
		Items: items,
	}
	return out, nil
}
