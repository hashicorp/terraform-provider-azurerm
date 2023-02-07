package cognitiveservicesaccounts

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

type AccountsListModelsOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AccountModel

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (AccountsListModelsOperationResponse, error)
}

type AccountsListModelsCompleteResult struct {
	Items []AccountModel
}

func (r AccountsListModelsOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r AccountsListModelsOperationResponse) LoadMore(ctx context.Context) (resp AccountsListModelsOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// AccountsListModels ...
func (c CognitiveServicesAccountsClient) AccountsListModels(ctx context.Context, id AccountId) (resp AccountsListModelsOperationResponse, err error) {
	req, err := c.preparerForAccountsListModels(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListModels", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListModels", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForAccountsListModels(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListModels", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForAccountsListModels prepares the AccountsListModels request.
func (c CognitiveServicesAccountsClient) preparerForAccountsListModels(ctx context.Context, id AccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/models", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForAccountsListModelsWithNextLink prepares the AccountsListModels request with the given nextLink token.
func (c CognitiveServicesAccountsClient) preparerForAccountsListModelsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForAccountsListModels handles the response to the AccountsListModels request. The method always
// closes the http.Response Body.
func (c CognitiveServicesAccountsClient) responderForAccountsListModels(resp *http.Response) (result AccountsListModelsOperationResponse, err error) {
	type page struct {
		Values   []AccountModel `json:"value"`
		NextLink *string        `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result AccountsListModelsOperationResponse, err error) {
			req, err := c.preparerForAccountsListModelsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListModels", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListModels", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForAccountsListModels(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "AccountsListModels", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// AccountsListModelsComplete retrieves all of the results into a single object
func (c CognitiveServicesAccountsClient) AccountsListModelsComplete(ctx context.Context, id AccountId) (AccountsListModelsCompleteResult, error) {
	return c.AccountsListModelsCompleteMatchingPredicate(ctx, id, AccountModelOperationPredicate{})
}

// AccountsListModelsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c CognitiveServicesAccountsClient) AccountsListModelsCompleteMatchingPredicate(ctx context.Context, id AccountId, predicate AccountModelOperationPredicate) (resp AccountsListModelsCompleteResult, err error) {
	items := make([]AccountModel, 0)

	page, err := c.AccountsListModels(ctx, id)
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

	out := AccountsListModelsCompleteResult{
		Items: items,
	}
	return out, nil
}
