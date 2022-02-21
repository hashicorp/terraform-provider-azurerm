package collection

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type ServicesListResponse struct {
	HttpResponse *http.Response
	Model        *[]ServicesDescription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ServicesListResponse, error)
}

type ServicesListCompleteResult struct {
	Items []ServicesDescription
}

func (r ServicesListResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ServicesListResponse) LoadMore(ctx context.Context) (resp ServicesListResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ServicesList ...
func (c CollectionClient) ServicesList(ctx context.Context, id commonids.SubscriptionId) (resp ServicesListResponse, err error) {
	req, err := c.preparerForServicesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "collection.CollectionClient", "ServicesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "collection.CollectionClient", "ServicesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForServicesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "collection.CollectionClient", "ServicesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ServicesListComplete retrieves all of the results into a single object
func (c CollectionClient) ServicesListComplete(ctx context.Context, id commonids.SubscriptionId) (ServicesListCompleteResult, error) {
	return c.ServicesListCompleteMatchingPredicate(ctx, id, ServicesDescriptionPredicate{})
}

// ServicesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c CollectionClient) ServicesListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ServicesDescriptionPredicate) (resp ServicesListCompleteResult, err error) {
	items := make([]ServicesDescription, 0)

	page, err := c.ServicesList(ctx, id)
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

	out := ServicesListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForServicesList prepares the ServicesList request.
func (c CollectionClient) preparerForServicesList(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.HealthcareApis/services", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForServicesListWithNextLink prepares the ServicesList request with the given nextLink token.
func (c CollectionClient) preparerForServicesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForServicesList handles the response to the ServicesList request. The method always
// closes the http.Response Body.
func (c CollectionClient) responderForServicesList(resp *http.Response) (result ServicesListResponse, err error) {
	type page struct {
		Values   []ServicesDescription `json:"value"`
		NextLink *string               `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ServicesListResponse, err error) {
			req, err := c.preparerForServicesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "collection.CollectionClient", "ServicesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "collection.CollectionClient", "ServicesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForServicesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "collection.CollectionClient", "ServicesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
