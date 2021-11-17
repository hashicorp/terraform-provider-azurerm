package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type ListByDomainTopicResponse struct {
	HttpResponse *http.Response
	Model        *[]EventSubscription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByDomainTopicResponse, error)
}

type ListByDomainTopicCompleteResult struct {
	Items []EventSubscription
}

func (r ListByDomainTopicResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByDomainTopicResponse) LoadMore(ctx context.Context) (resp ListByDomainTopicResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByDomainTopicOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListByDomainTopicOptions() ListByDomainTopicOptions {
	return ListByDomainTopicOptions{}
}

func (o ListByDomainTopicOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// ListByDomainTopic ...
func (c EventSubscriptionsClient) ListByDomainTopic(ctx context.Context, id DomainTopicId, options ListByDomainTopicOptions) (resp ListByDomainTopicResponse, err error) {
	req, err := c.preparerForListByDomainTopic(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListByDomainTopic", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListByDomainTopic", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByDomainTopic(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListByDomainTopic", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListByDomainTopicComplete retrieves all of the results into a single object
func (c EventSubscriptionsClient) ListByDomainTopicComplete(ctx context.Context, id DomainTopicId, options ListByDomainTopicOptions) (ListByDomainTopicCompleteResult, error) {
	return c.ListByDomainTopicCompleteMatchingPredicate(ctx, id, options, EventSubscriptionPredicate{})
}

// ListByDomainTopicCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventSubscriptionsClient) ListByDomainTopicCompleteMatchingPredicate(ctx context.Context, id DomainTopicId, options ListByDomainTopicOptions, predicate EventSubscriptionPredicate) (resp ListByDomainTopicCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	page, err := c.ListByDomainTopic(ctx, id, options)
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

	out := ListByDomainTopicCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListByDomainTopic prepares the ListByDomainTopic request.
func (c EventSubscriptionsClient) preparerForListByDomainTopic(ctx context.Context, id DomainTopicId, options ListByDomainTopicOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.EventGrid/eventSubscriptions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByDomainTopicWithNextLink prepares the ListByDomainTopic request with the given nextLink token.
func (c EventSubscriptionsClient) preparerForListByDomainTopicWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByDomainTopic handles the response to the ListByDomainTopic request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForListByDomainTopic(resp *http.Response) (result ListByDomainTopicResponse, err error) {
	type page struct {
		Values   []EventSubscription `json:"value"`
		NextLink *string             `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByDomainTopicResponse, err error) {
			req, err := c.preparerForListByDomainTopicWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListByDomainTopic", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListByDomainTopic", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByDomainTopic(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "ListByDomainTopic", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
