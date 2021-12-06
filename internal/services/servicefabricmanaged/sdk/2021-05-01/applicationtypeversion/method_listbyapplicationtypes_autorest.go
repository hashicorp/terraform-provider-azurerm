package applicationtypeversion

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByApplicationTypesResponse struct {
	HttpResponse *http.Response
	Model        *[]ApplicationTypeVersionResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByApplicationTypesResponse, error)
}

type ListByApplicationTypesCompleteResult struct {
	Items []ApplicationTypeVersionResource
}

func (r ListByApplicationTypesResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByApplicationTypesResponse) LoadMore(ctx context.Context) (resp ListByApplicationTypesResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByApplicationTypes ...
func (c ApplicationTypeVersionClient) ListByApplicationTypes(ctx context.Context, id ApplicationTypeId) (resp ListByApplicationTypesResponse, err error) {
	req, err := c.preparerForListByApplicationTypes(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationtypeversion.ApplicationTypeVersionClient", "ListByApplicationTypes", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationtypeversion.ApplicationTypeVersionClient", "ListByApplicationTypes", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByApplicationTypes(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "applicationtypeversion.ApplicationTypeVersionClient", "ListByApplicationTypes", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByApplicationTypesComplete retrieves all of the results into a single object
func (c ApplicationTypeVersionClient) ListByApplicationTypesComplete(ctx context.Context, id ApplicationTypeId) (ListByApplicationTypesCompleteResult, error) {
	return c.ListByApplicationTypesCompleteMatchingPredicate(ctx, id, ApplicationTypeVersionResourcePredicate{})
}

// ListByApplicationTypesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ApplicationTypeVersionClient) ListByApplicationTypesCompleteMatchingPredicate(ctx context.Context, id ApplicationTypeId, predicate ApplicationTypeVersionResourcePredicate) (resp ListByApplicationTypesCompleteResult, err error) {
	items := make([]ApplicationTypeVersionResource, 0)

	page, err := c.ListByApplicationTypes(ctx, id)
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

	out := ListByApplicationTypesCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByApplicationTypes prepares the ListByApplicationTypes request.
func (c ApplicationTypeVersionClient) preparerForListByApplicationTypes(ctx context.Context, id ApplicationTypeId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/versions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByApplicationTypesWithNextLink prepares the ListByApplicationTypes request with the given nextLink token.
func (c ApplicationTypeVersionClient) preparerForListByApplicationTypesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByApplicationTypes handles the response to the ListByApplicationTypes request. The method always
// closes the http.Response Body.
func (c ApplicationTypeVersionClient) responderForListByApplicationTypes(resp *http.Response) (result ListByApplicationTypesResponse, err error) {
	type page struct {
		Values   []ApplicationTypeVersionResource `json:"value"`
		NextLink *string                          `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByApplicationTypesResponse, err error) {
			req, err := c.preparerForListByApplicationTypesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "applicationtypeversion.ApplicationTypeVersionClient", "ListByApplicationTypes", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "applicationtypeversion.ApplicationTypeVersionClient", "ListByApplicationTypes", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByApplicationTypes(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "applicationtypeversion.ApplicationTypeVersionClient", "ListByApplicationTypes", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
