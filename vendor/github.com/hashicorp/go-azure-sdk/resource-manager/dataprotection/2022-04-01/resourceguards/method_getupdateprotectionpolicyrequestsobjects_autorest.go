package resourceguards

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

type GetUpdateProtectionPolicyRequestsObjectsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DppBaseResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (GetUpdateProtectionPolicyRequestsObjectsOperationResponse, error)
}

type GetUpdateProtectionPolicyRequestsObjectsCompleteResult struct {
	Items []DppBaseResource
}

func (r GetUpdateProtectionPolicyRequestsObjectsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r GetUpdateProtectionPolicyRequestsObjectsOperationResponse) LoadMore(ctx context.Context) (resp GetUpdateProtectionPolicyRequestsObjectsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// GetUpdateProtectionPolicyRequestsObjects ...
func (c ResourceGuardsClient) GetUpdateProtectionPolicyRequestsObjects(ctx context.Context, id ResourceGuardId) (resp GetUpdateProtectionPolicyRequestsObjectsOperationResponse, err error) {
	req, err := c.preparerForGetUpdateProtectionPolicyRequestsObjects(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetUpdateProtectionPolicyRequestsObjects", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetUpdateProtectionPolicyRequestsObjects", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForGetUpdateProtectionPolicyRequestsObjects(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetUpdateProtectionPolicyRequestsObjects", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// GetUpdateProtectionPolicyRequestsObjectsComplete retrieves all of the results into a single object
func (c ResourceGuardsClient) GetUpdateProtectionPolicyRequestsObjectsComplete(ctx context.Context, id ResourceGuardId) (GetUpdateProtectionPolicyRequestsObjectsCompleteResult, error) {
	return c.GetUpdateProtectionPolicyRequestsObjectsCompleteMatchingPredicate(ctx, id, DppBaseResourceOperationPredicate{})
}

// GetUpdateProtectionPolicyRequestsObjectsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ResourceGuardsClient) GetUpdateProtectionPolicyRequestsObjectsCompleteMatchingPredicate(ctx context.Context, id ResourceGuardId, predicate DppBaseResourceOperationPredicate) (resp GetUpdateProtectionPolicyRequestsObjectsCompleteResult, err error) {
	items := make([]DppBaseResource, 0)

	page, err := c.GetUpdateProtectionPolicyRequestsObjects(ctx, id)
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

	out := GetUpdateProtectionPolicyRequestsObjectsCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForGetUpdateProtectionPolicyRequestsObjects prepares the GetUpdateProtectionPolicyRequestsObjects request.
func (c ResourceGuardsClient) preparerForGetUpdateProtectionPolicyRequestsObjects(ctx context.Context, id ResourceGuardId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/updateProtectionPolicyRequests", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForGetUpdateProtectionPolicyRequestsObjectsWithNextLink prepares the GetUpdateProtectionPolicyRequestsObjects request with the given nextLink token.
func (c ResourceGuardsClient) preparerForGetUpdateProtectionPolicyRequestsObjectsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForGetUpdateProtectionPolicyRequestsObjects handles the response to the GetUpdateProtectionPolicyRequestsObjects request. The method always
// closes the http.Response Body.
func (c ResourceGuardsClient) responderForGetUpdateProtectionPolicyRequestsObjects(resp *http.Response) (result GetUpdateProtectionPolicyRequestsObjectsOperationResponse, err error) {
	type page struct {
		Values   []DppBaseResource `json:"value"`
		NextLink *string           `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result GetUpdateProtectionPolicyRequestsObjectsOperationResponse, err error) {
			req, err := c.preparerForGetUpdateProtectionPolicyRequestsObjectsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetUpdateProtectionPolicyRequestsObjects", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetUpdateProtectionPolicyRequestsObjects", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForGetUpdateProtectionPolicyRequestsObjects(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "resourceguards.ResourceGuardsClient", "GetUpdateProtectionPolicyRequestsObjects", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
