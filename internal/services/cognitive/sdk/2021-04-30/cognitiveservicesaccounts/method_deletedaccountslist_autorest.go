package cognitiveservicesaccounts

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DeletedAccountsListResponse struct {
	HttpResponse *http.Response
	Model        *[]Account

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (DeletedAccountsListResponse, error)
}

type DeletedAccountsListCompleteResult struct {
	Items []Account
}

func (r DeletedAccountsListResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r DeletedAccountsListResponse) LoadMore(ctx context.Context) (resp DeletedAccountsListResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// DeletedAccountsList ...
func (c CognitiveServicesAccountsClient) DeletedAccountsList(ctx context.Context, id SubscriptionId) (resp DeletedAccountsListResponse, err error) {
	req, err := c.preparerForDeletedAccountsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "DeletedAccountsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "DeletedAccountsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForDeletedAccountsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "DeletedAccountsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// DeletedAccountsListComplete retrieves all of the results into a single object
func (c CognitiveServicesAccountsClient) DeletedAccountsListComplete(ctx context.Context, id SubscriptionId) (DeletedAccountsListCompleteResult, error) {
	return c.DeletedAccountsListCompleteMatchingPredicate(ctx, id, AccountPredicate{})
}

// DeletedAccountsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c CognitiveServicesAccountsClient) DeletedAccountsListCompleteMatchingPredicate(ctx context.Context, id SubscriptionId, predicate AccountPredicate) (resp DeletedAccountsListCompleteResult, err error) {
	items := make([]Account, 0)

	page, err := c.DeletedAccountsList(ctx, id)
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

	out := DeletedAccountsListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForDeletedAccountsList prepares the DeletedAccountsList request.
func (c CognitiveServicesAccountsClient) preparerForDeletedAccountsList(ctx context.Context, id SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.CognitiveServices/deletedAccounts", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForDeletedAccountsListWithNextLink prepares the DeletedAccountsList request with the given nextLink token.
func (c CognitiveServicesAccountsClient) preparerForDeletedAccountsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForDeletedAccountsList handles the response to the DeletedAccountsList request. The method always
// closes the http.Response Body.
func (c CognitiveServicesAccountsClient) responderForDeletedAccountsList(resp *http.Response) (result DeletedAccountsListResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result DeletedAccountsListResponse, err error) {
			req, err := c.preparerForDeletedAccountsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "DeletedAccountsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "DeletedAccountsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForDeletedAccountsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "DeletedAccountsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
