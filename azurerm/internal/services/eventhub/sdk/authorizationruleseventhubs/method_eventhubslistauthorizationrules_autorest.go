package authorizationruleseventhubs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type EventHubsListAuthorizationRulesResponse struct {
	HttpResponse *http.Response
	Model        *[]AuthorizationRule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (EventHubsListAuthorizationRulesResponse, error)
}

type EventHubsListAuthorizationRulesCompleteResult struct {
	Items []AuthorizationRule
}

func (r EventHubsListAuthorizationRulesResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r EventHubsListAuthorizationRulesResponse) LoadMore(ctx context.Context) (resp EventHubsListAuthorizationRulesResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type AuthorizationRulePredicate struct {
	// TODO: implement me
}

func (p AuthorizationRulePredicate) Matches(input AuthorizationRule) bool {
	// TODO: implement me
	// if p.Name != nil && input.Name != *p.Name {
	// 	return false
	// }

	return true
}

// EventHubsListAuthorizationRules ...
func (c AuthorizationRulesEventHubsClient) EventHubsListAuthorizationRules(ctx context.Context, id EventhubId) (resp EventHubsListAuthorizationRulesResponse, err error) {
	req, err := c.preparerForEventHubsListAuthorizationRules(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsListAuthorizationRules", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsListAuthorizationRules", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForEventHubsListAuthorizationRules(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsListAuthorizationRules", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// EventHubsListAuthorizationRulesCompleteMatchingPredicate retrieves all of the results into a single object
func (c AuthorizationRulesEventHubsClient) EventHubsListAuthorizationRulesComplete(ctx context.Context, id EventhubId) (EventHubsListAuthorizationRulesCompleteResult, error) {
	return c.EventHubsListAuthorizationRulesCompleteMatchingPredicate(ctx, id, AuthorizationRulePredicate{})
}

// EventHubsListAuthorizationRulesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AuthorizationRulesEventHubsClient) EventHubsListAuthorizationRulesCompleteMatchingPredicate(ctx context.Context, id EventhubId, predicate AuthorizationRulePredicate) (resp EventHubsListAuthorizationRulesCompleteResult, err error) {
	items := make([]AuthorizationRule, 0)

	page, err := c.EventHubsListAuthorizationRules(ctx, id)
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

	out := EventHubsListAuthorizationRulesCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForEventHubsListAuthorizationRules prepares the EventHubsListAuthorizationRules request.
func (c AuthorizationRulesEventHubsClient) preparerForEventHubsListAuthorizationRules(ctx context.Context, id EventhubId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/authorizationRules", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForEventHubsListAuthorizationRulesWithNextLink prepares the EventHubsListAuthorizationRules request with the given nextLink token.
func (c AuthorizationRulesEventHubsClient) preparerForEventHubsListAuthorizationRulesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForEventHubsListAuthorizationRules handles the response to the EventHubsListAuthorizationRules request. The method always
// closes the http.Response Body.
func (c AuthorizationRulesEventHubsClient) responderForEventHubsListAuthorizationRules(resp *http.Response) (result EventHubsListAuthorizationRulesResponse, err error) {
	type page struct {
		Values   []AuthorizationRule `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result EventHubsListAuthorizationRulesResponse, err error) {
			req, err := c.preparerForEventHubsListAuthorizationRulesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsListAuthorizationRules", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsListAuthorizationRules", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForEventHubsListAuthorizationRules(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "authorizationruleseventhubs.AuthorizationRulesEventHubsClient", "EventHubsListAuthorizationRules", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
