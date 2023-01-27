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

type ListByRuleOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DataCollectionRuleAssociationProxyOnlyResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByRuleOperationResponse, error)
}

type ListByRuleCompleteResult struct {
	Items []DataCollectionRuleAssociationProxyOnlyResource
}

func (r ListByRuleOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByRuleOperationResponse) LoadMore(ctx context.Context) (resp ListByRuleOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByRule ...
func (c DataCollectionRuleAssociationsClient) ListByRule(ctx context.Context, id DataCollectionRuleId) (resp ListByRuleOperationResponse, err error) {
	req, err := c.preparerForListByRule(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datacollectionruleassociations.DataCollectionRuleAssociationsClient", "ListByRule", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "datacollectionruleassociations.DataCollectionRuleAssociationsClient", "ListByRule", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByRule(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "datacollectionruleassociations.DataCollectionRuleAssociationsClient", "ListByRule", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByRule prepares the ListByRule request.
func (c DataCollectionRuleAssociationsClient) preparerForListByRule(ctx context.Context, id DataCollectionRuleId) (*http.Request, error) {
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

// preparerForListByRuleWithNextLink prepares the ListByRule request with the given nextLink token.
func (c DataCollectionRuleAssociationsClient) preparerForListByRuleWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByRule handles the response to the ListByRule request. The method always
// closes the http.Response Body.
func (c DataCollectionRuleAssociationsClient) responderForListByRule(resp *http.Response) (result ListByRuleOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByRuleOperationResponse, err error) {
			req, err := c.preparerForListByRuleWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "datacollectionruleassociations.DataCollectionRuleAssociationsClient", "ListByRule", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "datacollectionruleassociations.DataCollectionRuleAssociationsClient", "ListByRule", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByRule(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "datacollectionruleassociations.DataCollectionRuleAssociationsClient", "ListByRule", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByRuleComplete retrieves all of the results into a single object
func (c DataCollectionRuleAssociationsClient) ListByRuleComplete(ctx context.Context, id DataCollectionRuleId) (ListByRuleCompleteResult, error) {
	return c.ListByRuleCompleteMatchingPredicate(ctx, id, DataCollectionRuleAssociationProxyOnlyResourceOperationPredicate{})
}

// ListByRuleCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DataCollectionRuleAssociationsClient) ListByRuleCompleteMatchingPredicate(ctx context.Context, id DataCollectionRuleId, predicate DataCollectionRuleAssociationProxyOnlyResourceOperationPredicate) (resp ListByRuleCompleteResult, err error) {
	items := make([]DataCollectionRuleAssociationProxyOnlyResource, 0)

	page, err := c.ListByRule(ctx, id)
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

	out := ListByRuleCompleteResult{
		Items: items,
	}
	return out, nil
}
