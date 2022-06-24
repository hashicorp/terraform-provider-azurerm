package cognitiveservicesaccounts

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Account

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AccountsListByResourceGroupOperationResponse, error)
}

type AccountsListByResourceGroupCompleteResult struct {
	Items []Account
}

func (r AccountsListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AccountsListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp AccountsListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// AccountsListByResourceGroup ...
func (c CognitiveServicesAccountsClient) AccountsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp AccountsListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForAccountsListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAccountsListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// AccountsListByResourceGroupComplete retrieves all of the results into a single object
func (c CognitiveServicesAccountsClient) AccountsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (AccountsListByResourceGroupCompleteResult, error) {
	return c.AccountsListByResourceGroupCompleteMatchingPredicate(ctx, id, AccountOperationPredicate{})
}

// AccountsListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c CognitiveServicesAccountsClient) AccountsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate AccountOperationPredicate) (resp AccountsListByResourceGroupCompleteResult, err error) {
	items := make([]Account, 0)

	page, err := c.AccountsListByResourceGroup(ctx, id)
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

	out := AccountsListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForAccountsListByResourceGroup prepares the AccountsListByResourceGroup request.
func (c CognitiveServicesAccountsClient) preparerForAccountsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.CognitiveServices/accounts", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForAccountsListByResourceGroupWithNextLink prepares the AccountsListByResourceGroup request with the given nextLink token.
func (c CognitiveServicesAccountsClient) preparerForAccountsListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAccountsListByResourceGroup handles the response to the AccountsListByResourceGroup request. The method always
// closes the http.Response Body.
func (c CognitiveServicesAccountsClient) responderForAccountsListByResourceGroup(resp *http.Response) (result AccountsListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []Account `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AccountsListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForAccountsListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAccountsListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
