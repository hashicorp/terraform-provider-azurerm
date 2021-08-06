package eventhubs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByNamespaceResponse struct {
	HttpResponse *http.Response
	Model        *[]Eventhub

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByNamespaceResponse, error)
}

type ListByNamespaceCompleteResult struct {
	Items []Eventhub
}

func (r ListByNamespaceResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByNamespaceResponse) LoadMore(ctx context.Context) (resp ListByNamespaceResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByNamespaceOptions struct {
	Skip *int64
	Top  *int64
}

func DefaultListByNamespaceOptions() ListByNamespaceOptions {
	return ListByNamespaceOptions{}
}

func (o ListByNamespaceOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Skip != nil {
		out["$skip"] = *o.Skip
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

type EventhubPredicate struct {
	// TODO: implement me
}

func (p EventhubPredicate) Matches(input Eventhub) bool {
	// TODO: implement me
	// if p.Name != nil && input.Name != *p.Name {
	// 	return false
	// }

	return true
}

// ListByNamespace ...
func (c EventHubsClient) ListByNamespace(ctx context.Context, id NamespaceId, options ListByNamespaceOptions) (resp ListByNamespaceResponse, err error) {
	req, err := c.preparerForListByNamespace(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubs.EventHubsClient", "ListByNamespace", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubs.EventHubsClient", "ListByNamespace", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByNamespace(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventhubs.EventHubsClient", "ListByNamespace", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByNamespaceCompleteMatchingPredicate retrieves all of the results into a single object
func (c EventHubsClient) ListByNamespaceComplete(ctx context.Context, id NamespaceId, options ListByNamespaceOptions) (ListByNamespaceCompleteResult, error) {
	return c.ListByNamespaceCompleteMatchingPredicate(ctx, id, options, EventhubPredicate{})
}

// ListByNamespaceCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventHubsClient) ListByNamespaceCompleteMatchingPredicate(ctx context.Context, id NamespaceId, options ListByNamespaceOptions, predicate EventhubPredicate) (resp ListByNamespaceCompleteResult, err error) {
	items := make([]Eventhub, 0)

	page, err := c.ListByNamespace(ctx, id, options)
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

	out := ListByNamespaceCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByNamespace prepares the ListByNamespace request.
func (c EventHubsClient) preparerForListByNamespace(ctx context.Context, id NamespaceId, options ListByNamespaceOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/eventhubs", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByNamespaceWithNextLink prepares the ListByNamespace request with the given nextLink token.
func (c EventHubsClient) preparerForListByNamespaceWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByNamespace handles the response to the ListByNamespace request. The method always
// closes the http.Response Body.
func (c EventHubsClient) responderForListByNamespace(resp *http.Response) (result ListByNamespaceResponse, err error) {
	type page struct {
		Values   []Eventhub `json:"value"`
		NextLink *string    `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByNamespaceResponse, err error) {
			req, err := c.preparerForListByNamespaceWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventhubs.EventHubsClient", "ListByNamespace", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventhubs.EventHubsClient", "ListByNamespace", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByNamespace(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventhubs.EventHubsClient", "ListByNamespace", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
