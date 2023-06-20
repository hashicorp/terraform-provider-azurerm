package datacollectionruleassociations

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

type ListByDataCollectionEndpointOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DataCollectionRuleAssociationProxyOnlyResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByDataCollectionEndpointOperationResponse, error)
}

type ListByDataCollectionEndpointCompleteResult struct {
	Items []DataCollectionRuleAssociationProxyOnlyResource
}

func (r ListByDataCollectionEndpointOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByDataCollectionEndpointOperationResponse) LoadMore(ctx context.Context) (resp ListByDataCollectionEndpointOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByDataCollectionEndpoint ...
func (c DataCollectionRuleAssociationsClient) ListByDataCollectionEndpoint(ctx context.Context, id DataCollectionEndpointId) (resp ListByDataCollectionEndpointOperationResponse, err error) {
	req, err := c.preparerForListByDataCollectionEndpoint(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datacollectionruleassociations.DataCollectionRuleAssociationsClient", "ListByDataCollectionEndpoint", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "datacollectionruleassociations.DataCollectionRuleAssociationsClient", "ListByDataCollectionEndpoint", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByDataCollectionEndpoint(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datacollectionruleassociations.DataCollectionRuleAssociationsClient", "ListByDataCollectionEndpoint", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByDataCollectionEndpoint prepares the ListByDataCollectionEndpoint request.
func (c DataCollectionRuleAssociationsClient) preparerForListByDataCollectionEndpoint(ctx context.Context, id DataCollectionEndpointId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/associations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByDataCollectionEndpointWithNextLink prepares the ListByDataCollectionEndpoint request with the given nextLink token.
func (c DataCollectionRuleAssociationsClient) preparerForListByDataCollectionEndpointWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByDataCollectionEndpoint handles the response to the ListByDataCollectionEndpoint request. The method always
// closes the http.Response Body.
func (c DataCollectionRuleAssociationsClient) responderForListByDataCollectionEndpoint(resp *http.Response) (result ListByDataCollectionEndpointOperationResponse, err error) {
	type page struct {
		Values   []DataCollectionRuleAssociationProxyOnlyResource `json:"value"`
		NextLink *string                                          `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByDataCollectionEndpointOperationResponse, err error) {
			req, err := c.preparerForListByDataCollectionEndpointWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "datacollectionruleassociations.DataCollectionRuleAssociationsClient", "ListByDataCollectionEndpoint", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "datacollectionruleassociations.DataCollectionRuleAssociationsClient", "ListByDataCollectionEndpoint", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByDataCollectionEndpoint(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "datacollectionruleassociations.DataCollectionRuleAssociationsClient", "ListByDataCollectionEndpoint", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByDataCollectionEndpointComplete retrieves all of the results into a single object
func (c DataCollectionRuleAssociationsClient) ListByDataCollectionEndpointComplete(ctx context.Context, id DataCollectionEndpointId) (ListByDataCollectionEndpointCompleteResult, error) {
	return c.ListByDataCollectionEndpointCompleteMatchingPredicate(ctx, id, DataCollectionRuleAssociationProxyOnlyResourceOperationPredicate{})
}

// ListByDataCollectionEndpointCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DataCollectionRuleAssociationsClient) ListByDataCollectionEndpointCompleteMatchingPredicate(ctx context.Context, id DataCollectionEndpointId, predicate DataCollectionRuleAssociationProxyOnlyResourceOperationPredicate) (resp ListByDataCollectionEndpointCompleteResult, err error) {
	items := make([]DataCollectionRuleAssociationProxyOnlyResource, 0)

	page, err := c.ListByDataCollectionEndpoint(ctx, id)
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

	out := ListByDataCollectionEndpointCompleteResult{
		Items: items,
	}
	return out, nil
}
