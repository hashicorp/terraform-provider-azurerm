package operationsinacluster

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type OperationStatusListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]OperationStatusResult

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (OperationStatusListOperationResponse, error)
}

type OperationStatusListCompleteResult struct {
	Items []OperationStatusResult
}

func (r OperationStatusListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r OperationStatusListOperationResponse) LoadMore(ctx context.Context) (resp OperationStatusListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// OperationStatusList ...
func (c OperationsInAClusterClient) OperationStatusList(ctx context.Context, id ProviderId) (resp OperationStatusListOperationResponse, err error) {
	req, err := c.preparerForOperationStatusList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationsinacluster.OperationsInAClusterClient", "OperationStatusList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationsinacluster.OperationsInAClusterClient", "OperationStatusList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForOperationStatusList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "operationsinacluster.OperationsInAClusterClient", "OperationStatusList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// OperationStatusListComplete retrieves all of the results into a single object
func (c OperationsInAClusterClient) OperationStatusListComplete(ctx context.Context, id ProviderId) (OperationStatusListCompleteResult, error) {
	return c.OperationStatusListCompleteMatchingPredicate(ctx, id, OperationStatusResultOperationPredicate{})
}

// OperationStatusListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c OperationsInAClusterClient) OperationStatusListCompleteMatchingPredicate(ctx context.Context, id ProviderId, predicate OperationStatusResultOperationPredicate) (resp OperationStatusListCompleteResult, err error) {
	items := make([]OperationStatusResult, 0)

	page, err := c.OperationStatusList(ctx, id)
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

	out := OperationStatusListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForOperationStatusList prepares the OperationStatusList request.
func (c OperationsInAClusterClient) preparerForOperationStatusList(ctx context.Context, id ProviderId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.KubernetesConfiguration/operations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForOperationStatusListWithNextLink prepares the OperationStatusList request with the given nextLink token.
func (c OperationsInAClusterClient) preparerForOperationStatusListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForOperationStatusList handles the response to the OperationStatusList request. The method always
// closes the http.Response Body.
func (c OperationsInAClusterClient) responderForOperationStatusList(resp *http.Response) (result OperationStatusListOperationResponse, err error) {
	type page struct {
		Values   []OperationStatusResult `json:"value"`
		NextLink *string                 `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result OperationStatusListOperationResponse, err error) {
			req, err := c.preparerForOperationStatusListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationsinacluster.OperationsInAClusterClient", "OperationStatusList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationsinacluster.OperationsInAClusterClient", "OperationStatusList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForOperationStatusList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "operationsinacluster.OperationsInAClusterClient", "OperationStatusList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
