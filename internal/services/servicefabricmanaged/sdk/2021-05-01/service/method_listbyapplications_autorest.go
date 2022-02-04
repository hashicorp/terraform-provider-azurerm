package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByApplicationsResponse struct {
	HttpResponse *http.Response
	Model        *[]ServiceResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByApplicationsResponse, error)
}

type ListByApplicationsCompleteResult struct {
	Items []ServiceResource
}

func (r ListByApplicationsResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByApplicationsResponse) LoadMore(ctx context.Context) (resp ListByApplicationsResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByApplications ...
func (c ServiceClient) ListByApplications(ctx context.Context, id ApplicationId) (resp ListByApplicationsResponse, err error) {
	req, err := c.preparerForListByApplications(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "service.ServiceClient", "ListByApplications", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "service.ServiceClient", "ListByApplications", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByApplications(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "service.ServiceClient", "ListByApplications", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByApplicationsComplete retrieves all of the results into a single object
func (c ServiceClient) ListByApplicationsComplete(ctx context.Context, id ApplicationId) (ListByApplicationsCompleteResult, error) {
	return c.ListByApplicationsCompleteMatchingPredicate(ctx, id, ServiceResourcePredicate{})
}

// ListByApplicationsCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ServiceClient) ListByApplicationsCompleteMatchingPredicate(ctx context.Context, id ApplicationId, predicate ServiceResourcePredicate) (resp ListByApplicationsCompleteResult, err error) {
	items := make([]ServiceResource, 0)

	page, err := c.ListByApplications(ctx, id)
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

	out := ListByApplicationsCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByApplications prepares the ListByApplications request.
func (c ServiceClient) preparerForListByApplications(ctx context.Context, id ApplicationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/services", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByApplicationsWithNextLink prepares the ListByApplications request with the given nextLink token.
func (c ServiceClient) preparerForListByApplicationsWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByApplications handles the response to the ListByApplications request. The method always
// closes the http.Response Body.
func (c ServiceClient) responderForListByApplications(resp *http.Response) (result ListByApplicationsResponse, err error) {
	type page struct {
		Values   []ServiceResource `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByApplicationsResponse, err error) {
			req, err := c.preparerForListByApplicationsWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "service.ServiceClient", "ListByApplications", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "service.ServiceClient", "ListByApplications", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByApplications(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "service.ServiceClient", "ListByApplications", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
