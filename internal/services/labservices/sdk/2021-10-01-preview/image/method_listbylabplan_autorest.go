package image

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByLabPlanResponse struct {
	HttpResponse *http.Response
	Model        *[]Image

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByLabPlanResponse, error)
}

type ListByLabPlanCompleteResult struct {
	Items []Image
}

func (r ListByLabPlanResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByLabPlanResponse) LoadMore(ctx context.Context) (resp ListByLabPlanResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByLabPlan ...
func (c ImageClient) ListByLabPlan(ctx context.Context, id LabPlanId) (resp ListByLabPlanResponse, err error) {
	req, err := c.preparerForListByLabPlan(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "image.ImageClient", "ListByLabPlan", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "image.ImageClient", "ListByLabPlan", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByLabPlan(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "image.ImageClient", "ListByLabPlan", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByLabPlanComplete retrieves all of the results into a single object
func (c ImageClient) ListByLabPlanComplete(ctx context.Context, id LabPlanId) (ListByLabPlanCompleteResult, error) {
	return c.ListByLabPlanCompleteMatchingPredicate(ctx, id, ImagePredicate{})
}

// ListByLabPlanCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ImageClient) ListByLabPlanCompleteMatchingPredicate(ctx context.Context, id LabPlanId, predicate ImagePredicate) (resp ListByLabPlanCompleteResult, err error) {
	items := make([]Image, 0)

	page, err := c.ListByLabPlan(ctx, id)
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

	out := ListByLabPlanCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByLabPlan prepares the ListByLabPlan request.
func (c ImageClient) preparerForListByLabPlan(ctx context.Context, id LabPlanId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/images", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByLabPlanWithNextLink prepares the ListByLabPlan request with the given nextLink token.
func (c ImageClient) preparerForListByLabPlanWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByLabPlan handles the response to the ListByLabPlan request. The method always
// closes the http.Response Body.
func (c ImageClient) responderForListByLabPlan(resp *http.Response) (result ListByLabPlanResponse, err error) {
	type page struct {
		Values   []Image `json:"value"`
		NextLink *string `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByLabPlanResponse, err error) {
			req, err := c.preparerForListByLabPlanWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "image.ImageClient", "ListByLabPlan", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "image.ImageClient", "ListByLabPlan", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByLabPlan(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "image.ImageClient", "ListByLabPlan", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
