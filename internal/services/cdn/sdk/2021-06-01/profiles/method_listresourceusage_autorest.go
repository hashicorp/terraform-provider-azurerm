package profiles

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListResourceUsageResponse struct {
	HttpResponse *http.Response
	Model        *[]ResourceUsage

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListResourceUsageResponse, error)
}

type ListResourceUsageCompleteResult struct {
	Items []ResourceUsage
}

func (r ListResourceUsageResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListResourceUsageResponse) LoadMore(ctx context.Context) (resp ListResourceUsageResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListResourceUsage ...
func (c ProfilesClient) ListResourceUsage(ctx context.Context, id ProfileId) (resp ListResourceUsageResponse, err error) {
	req, err := c.preparerForListResourceUsage(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "ListResourceUsage", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "ListResourceUsage", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListResourceUsage(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "ListResourceUsage", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListResourceUsageComplete retrieves all of the results into a single object
func (c ProfilesClient) ListResourceUsageComplete(ctx context.Context, id ProfileId) (ListResourceUsageCompleteResult, error) {
	return c.ListResourceUsageCompleteMatchingPredicate(ctx, id, ResourceUsagePredicate{})
}

// ListResourceUsageCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ProfilesClient) ListResourceUsageCompleteMatchingPredicate(ctx context.Context, id ProfileId, predicate ResourceUsagePredicate) (resp ListResourceUsageCompleteResult, err error) {
	items := make([]ResourceUsage, 0)

	page, err := c.ListResourceUsage(ctx, id)
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

	out := ListResourceUsageCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListResourceUsage prepares the ListResourceUsage request.
func (c ProfilesClient) preparerForListResourceUsage(ctx context.Context, id ProfileId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/checkResourceUsage", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListResourceUsageWithNextLink prepares the ListResourceUsage request with the given nextLink token.
func (c ProfilesClient) preparerForListResourceUsageWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListResourceUsage handles the response to the ListResourceUsage request. The method always
// closes the http.Response Body.
func (c ProfilesClient) responderForListResourceUsage(resp *http.Response) (result ListResourceUsageResponse, err error) {
	type page struct {
		Values   []ResourceUsage `json:"value"`
		NextLink *string         `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListResourceUsageResponse, err error) {
			req, err := c.preparerForListResourceUsageWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "ListResourceUsage", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "ListResourceUsage", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListResourceUsage(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "profiles.ProfilesClient", "ListResourceUsage", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
