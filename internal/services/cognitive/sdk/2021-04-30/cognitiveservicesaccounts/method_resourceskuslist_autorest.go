package cognitiveservicesaccounts

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ResourceSkusListResponse struct {
	HttpResponse *http.Response
	Model        *[]ResourceSku

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ResourceSkusListResponse, error)
}

type ResourceSkusListCompleteResult struct {
	Items []ResourceSku
}

func (r ResourceSkusListResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ResourceSkusListResponse) LoadMore(ctx context.Context) (resp ResourceSkusListResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ResourceSkusList ...
func (c CognitiveServicesAccountsClient) ResourceSkusList(ctx context.Context, id SubscriptionId) (resp ResourceSkusListResponse, err error) {
	req, err := c.preparerForResourceSkusList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "ResourceSkusList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "ResourceSkusList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForResourceSkusList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "ResourceSkusList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ResourceSkusListComplete retrieves all of the results into a single object
func (c CognitiveServicesAccountsClient) ResourceSkusListComplete(ctx context.Context, id SubscriptionId) (ResourceSkusListCompleteResult, error) {
	return c.ResourceSkusListCompleteMatchingPredicate(ctx, id, ResourceSkuPredicate{})
}

// ResourceSkusListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c CognitiveServicesAccountsClient) ResourceSkusListCompleteMatchingPredicate(ctx context.Context, id SubscriptionId, predicate ResourceSkuPredicate) (resp ResourceSkusListCompleteResult, err error) {
	items := make([]ResourceSku, 0)

	page, err := c.ResourceSkusList(ctx, id)
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

	out := ResourceSkusListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForResourceSkusList prepares the ResourceSkusList request.
func (c CognitiveServicesAccountsClient) preparerForResourceSkusList(ctx context.Context, id SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.CognitiveServices/skus", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForResourceSkusListWithNextLink prepares the ResourceSkusList request with the given nextLink token.
func (c CognitiveServicesAccountsClient) preparerForResourceSkusListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForResourceSkusList handles the response to the ResourceSkusList request. The method always
// closes the http.Response Body.
func (c CognitiveServicesAccountsClient) responderForResourceSkusList(resp *http.Response) (result ResourceSkusListResponse, err error) {
	type page struct {
		Values   []ResourceSku `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ResourceSkusListResponse, err error) {
			req, err := c.preparerForResourceSkusListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "ResourceSkusList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "ResourceSkusList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForResourceSkusList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "cognitiveservicesaccounts.CognitiveServicesAccountsClient", "ResourceSkusList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
