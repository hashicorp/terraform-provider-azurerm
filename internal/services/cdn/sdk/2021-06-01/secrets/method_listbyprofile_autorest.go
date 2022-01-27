package secrets

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByProfileResponse struct {
	HttpResponse *http.Response
	Model        *[]Secret

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByProfileResponse, error)
}

type ListByProfileCompleteResult struct {
	Items []Secret
}

func (r ListByProfileResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByProfileResponse) LoadMore(ctx context.Context) (resp ListByProfileResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListByProfile ...
func (c SecretsClient) ListByProfile(ctx context.Context, id ProfileId) (resp ListByProfileResponse, err error) {
	req, err := c.preparerForListByProfile(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "secrets.SecretsClient", "ListByProfile", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "secrets.SecretsClient", "ListByProfile", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByProfile(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "secrets.SecretsClient", "ListByProfile", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByProfileComplete retrieves all of the results into a single object
func (c SecretsClient) ListByProfileComplete(ctx context.Context, id ProfileId) (ListByProfileCompleteResult, error) {
	return c.ListByProfileCompleteMatchingPredicate(ctx, id, SecretPredicate{})
}

// ListByProfileCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SecretsClient) ListByProfileCompleteMatchingPredicate(ctx context.Context, id ProfileId, predicate SecretPredicate) (resp ListByProfileCompleteResult, err error) {
	items := make([]Secret, 0)

	page, err := c.ListByProfile(ctx, id)
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

	out := ListByProfileCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByProfile prepares the ListByProfile request.
func (c SecretsClient) preparerForListByProfile(ctx context.Context, id ProfileId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/secrets", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByProfileWithNextLink prepares the ListByProfile request with the given nextLink token.
func (c SecretsClient) preparerForListByProfileWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByProfile handles the response to the ListByProfile request. The method always
// closes the http.Response Body.
func (c SecretsClient) responderForListByProfile(resp *http.Response) (result ListByProfileResponse, err error) {
	type page struct {
		Values   []Secret `json:"value"`
		NextLink *string  `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByProfileResponse, err error) {
			req, err := c.preparerForListByProfileWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "secrets.SecretsClient", "ListByProfile", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "secrets.SecretsClient", "ListByProfile", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByProfile(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "secrets.SecretsClient", "ListByProfile", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
