package eventsubscriptions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type PartnerTopicEventSubscriptionsListByPartnerTopicResponse struct {
	HttpResponse *http.Response
	Model        *[]EventSubscription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (PartnerTopicEventSubscriptionsListByPartnerTopicResponse, error)
}

type PartnerTopicEventSubscriptionsListByPartnerTopicCompleteResult struct {
	Items []EventSubscription
}

func (r PartnerTopicEventSubscriptionsListByPartnerTopicResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r PartnerTopicEventSubscriptionsListByPartnerTopicResponse) LoadMore(ctx context.Context) (resp PartnerTopicEventSubscriptionsListByPartnerTopicResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type PartnerTopicEventSubscriptionsListByPartnerTopicOptions struct {
	Filter *string
	Top    *int64
}

func DefaultPartnerTopicEventSubscriptionsListByPartnerTopicOptions() PartnerTopicEventSubscriptionsListByPartnerTopicOptions {
	return PartnerTopicEventSubscriptionsListByPartnerTopicOptions{}
}

func (o PartnerTopicEventSubscriptionsListByPartnerTopicOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Top != nil {
		out["$top"] = *o.Top
	}

	return out
}

// PartnerTopicEventSubscriptionsListByPartnerTopic ...
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsListByPartnerTopic(ctx context.Context, id PartnerTopicId, options PartnerTopicEventSubscriptionsListByPartnerTopicOptions) (resp PartnerTopicEventSubscriptionsListByPartnerTopicResponse, err error) {
	req, err := c.preparerForPartnerTopicEventSubscriptionsListByPartnerTopic(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsListByPartnerTopic", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsListByPartnerTopic", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForPartnerTopicEventSubscriptionsListByPartnerTopic(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsListByPartnerTopic", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// PartnerTopicEventSubscriptionsListByPartnerTopicComplete retrieves all of the results into a single object
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsListByPartnerTopicComplete(ctx context.Context, id PartnerTopicId, options PartnerTopicEventSubscriptionsListByPartnerTopicOptions) (PartnerTopicEventSubscriptionsListByPartnerTopicCompleteResult, error) {
	return c.PartnerTopicEventSubscriptionsListByPartnerTopicCompleteMatchingPredicate(ctx, id, options, EventSubscriptionPredicate{})
}

// PartnerTopicEventSubscriptionsListByPartnerTopicCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c EventSubscriptionsClient) PartnerTopicEventSubscriptionsListByPartnerTopicCompleteMatchingPredicate(ctx context.Context, id PartnerTopicId, options PartnerTopicEventSubscriptionsListByPartnerTopicOptions, predicate EventSubscriptionPredicate) (resp PartnerTopicEventSubscriptionsListByPartnerTopicCompleteResult, err error) {
	items := make([]EventSubscription, 0)

	page, err := c.PartnerTopicEventSubscriptionsListByPartnerTopic(ctx, id, options)
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

	out := PartnerTopicEventSubscriptionsListByPartnerTopicCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForPartnerTopicEventSubscriptionsListByPartnerTopic prepares the PartnerTopicEventSubscriptionsListByPartnerTopic request.
func (c EventSubscriptionsClient) preparerForPartnerTopicEventSubscriptionsListByPartnerTopic(ctx context.Context, id PartnerTopicId, options PartnerTopicEventSubscriptionsListByPartnerTopicOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/eventSubscriptions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForPartnerTopicEventSubscriptionsListByPartnerTopicWithNextLink prepares the PartnerTopicEventSubscriptionsListByPartnerTopic request with the given nextLink token.
func (c EventSubscriptionsClient) preparerForPartnerTopicEventSubscriptionsListByPartnerTopicWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForPartnerTopicEventSubscriptionsListByPartnerTopic handles the response to the PartnerTopicEventSubscriptionsListByPartnerTopic request. The method always
// closes the http.Response Body.
func (c EventSubscriptionsClient) responderForPartnerTopicEventSubscriptionsListByPartnerTopic(resp *http.Response) (result PartnerTopicEventSubscriptionsListByPartnerTopicResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result PartnerTopicEventSubscriptionsListByPartnerTopicResponse, err error) {
			req, err := c.preparerForPartnerTopicEventSubscriptionsListByPartnerTopicWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsListByPartnerTopic", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsListByPartnerTopic", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForPartnerTopicEventSubscriptionsListByPartnerTopic(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "eventsubscriptions.EventSubscriptionsClient", "PartnerTopicEventSubscriptionsListByPartnerTopic", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
