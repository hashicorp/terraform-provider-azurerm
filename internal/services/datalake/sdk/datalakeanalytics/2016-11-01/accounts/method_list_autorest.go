package accounts

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListResponse struct {
	HttpResponse *http.Response
	Model        *[]DataLakeAnalyticsAccountBasic

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListResponse, error)
}

type ListCompleteResult struct {
	Items []DataLakeAnalyticsAccountBasic
}

func (r ListResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListResponse) LoadMore(ctx context.Context) (resp ListResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListOptions struct {
	Count   *bool
	Filter  *string
	Orderby *string
	Select  *string
	Skip    *int64
	Top     *int64
}

func DefaultListOptions() ListOptions {
	return ListOptions{}
}

func (o ListOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Count != nil {
		out["$count"] = *o.Count
	}

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Orderby != nil {
		out["$orderby"] = *o.Orderby
	}

	if o.Select != nil {
		out["$select"] = *o.Select
	}

	if o.Skip != nil {
		out["$skip"] = *o.Skip
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// List ...
func (c AccountsClient) List(ctx context.Context, id SubscriptionId, options ListOptions) (resp ListResponse, err error) {
	req, err := c.preparerForList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "List", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "List", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "List", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListComplete retrieves all of the results into a single object
func (c AccountsClient) ListComplete(ctx context.Context, id SubscriptionId, options ListOptions) (ListCompleteResult, error) {
	return c.ListCompleteMatchingPredicate(ctx, id, options, DataLakeAnalyticsAccountBasicPredicate{})
}

// ListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AccountsClient) ListCompleteMatchingPredicate(ctx context.Context, id SubscriptionId, options ListOptions, predicate DataLakeAnalyticsAccountBasicPredicate) (resp ListCompleteResult, err error) {
	items := make([]DataLakeAnalyticsAccountBasic, 0)

	page, err := c.List(ctx, id, options)
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

	out := ListCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForList prepares the List request.
func (c AccountsClient) preparerForList(ctx context.Context, id SubscriptionId, options ListOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.DataLakeAnalytics/accounts", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListWithNextLink prepares the List request with the given nextLink token.
func (c AccountsClient) preparerForListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForList handles the response to the List request. The method always
// closes the http.Response Body.
func (c AccountsClient) responderForList(resp *http.Response) (result ListResponse, err error) {
	type page struct {
		Values   []DataLakeAnalyticsAccountBasic `json:"value"`
		NextLink *string                         `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListResponse, err error) {
			req, err := c.preparerForListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "List", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "List", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "accounts.AccountsClient", "List", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
