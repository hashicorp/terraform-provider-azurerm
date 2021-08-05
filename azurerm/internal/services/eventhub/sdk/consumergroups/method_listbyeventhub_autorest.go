package consumergroups

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByEventHubResponse struct {
	HttpResponse *http.Response
	Model        *[]ConsumerGroup

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByEventHubResponse, error)
}

type ListByEventHubCompleteResult struct {
	Items []ConsumerGroup
}

func (r ListByEventHubResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByEventHubResponse) LoadMore(ctx context.Context) (resp ListByEventHubResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByEventHubOptions struct {
	Skip *int64
	Top  *int64
}

func DefaultListByEventHubOptions() ListByEventHubOptions {
	return ListByEventHubOptions{}
}

func (o ListByEventHubOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Skip != nil {
		out["$skip"] = *o.Skip
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

type ConsumerGroupPredicate struct {
	// TODO: implement me
}

func (p ConsumerGroupPredicate) Matches(input ConsumerGroup) bool {
	// TODO: implement me
	// if p.Name != nil && input.Name != *p.Name {
	// 	return false
	// }

	return true
}

// ListByEventHub ...
func (c ConsumerGroupsClient) ListByEventHub(ctx context.Context, id EventhubId, options ListByEventHubOptions) (resp ListByEventHubResponse, err error) {
	req, err := c.preparerForListByEventHub(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumergroups.ConsumerGroupsClient", "ListByEventHub", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumergroups.ConsumerGroupsClient", "ListByEventHub", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByEventHub(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "consumergroups.ConsumerGroupsClient", "ListByEventHub", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByEventHubCompleteMatchingPredicate retrieves all of the results into a single object
func (c ConsumerGroupsClient) ListByEventHubComplete(ctx context.Context, id EventhubId, options ListByEventHubOptions) (ListByEventHubCompleteResult, error) {
	return c.ListByEventHubCompleteMatchingPredicate(ctx, id, options, ConsumerGroupPredicate{})
}

// ListByEventHubCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ConsumerGroupsClient) ListByEventHubCompleteMatchingPredicate(ctx context.Context, id EventhubId, options ListByEventHubOptions, predicate ConsumerGroupPredicate) (resp ListByEventHubCompleteResult, err error) {
	items := make([]ConsumerGroup, 0)

	page, err := c.ListByEventHub(ctx, id, options)
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

	out := ListByEventHubCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByEventHub prepares the ListByEventHub request.
func (c ConsumerGroupsClient) preparerForListByEventHub(ctx context.Context, id EventhubId, options ListByEventHubOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/consumergroups", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByEventHubWithNextLink prepares the ListByEventHub request with the given nextLink token.
func (c ConsumerGroupsClient) preparerForListByEventHubWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByEventHub handles the response to the ListByEventHub request. The method always
// closes the http.Response Body.
func (c ConsumerGroupsClient) responderForListByEventHub(resp *http.Response) (result ListByEventHubResponse, err error) {
	type page struct {
		Values   []ConsumerGroup `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByEventHubResponse, err error) {
			req, err := c.preparerForListByEventHubWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "consumergroups.ConsumerGroupsClient", "ListByEventHub", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "consumergroups.ConsumerGroupsClient", "ListByEventHub", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByEventHub(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "consumergroups.ConsumerGroupsClient", "ListByEventHub", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
