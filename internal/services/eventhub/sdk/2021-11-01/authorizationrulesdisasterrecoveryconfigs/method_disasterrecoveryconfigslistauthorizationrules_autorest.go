package authorizationrulesdisasterrecoveryconfigs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type DisasterRecoveryConfigsListAuthorizationRulesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AuthorizationRule

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (DisasterRecoveryConfigsListAuthorizationRulesOperationResponse, error)
}

type DisasterRecoveryConfigsListAuthorizationRulesCompleteResult struct {
	Items []AuthorizationRule
}

func (r DisasterRecoveryConfigsListAuthorizationRulesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r DisasterRecoveryConfigsListAuthorizationRulesOperationResponse) LoadMore(ctx context.Context) (resp DisasterRecoveryConfigsListAuthorizationRulesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// DisasterRecoveryConfigsListAuthorizationRules ...
func (c AuthorizationRulesDisasterRecoveryConfigsClient) DisasterRecoveryConfigsListAuthorizationRules(ctx context.Context, id DisasterRecoveryConfigId) (resp DisasterRecoveryConfigsListAuthorizationRulesOperationResponse, err error) {
	req, err := c.preparerForDisasterRecoveryConfigsListAuthorizationRules(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesdisasterrecoveryconfigs.AuthorizationRulesDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsListAuthorizationRules", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesdisasterrecoveryconfigs.AuthorizationRulesDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsListAuthorizationRules", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForDisasterRecoveryConfigsListAuthorizationRules(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorizationrulesdisasterrecoveryconfigs.AuthorizationRulesDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsListAuthorizationRules", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// DisasterRecoveryConfigsListAuthorizationRulesComplete retrieves all of the results into a single object
func (c AuthorizationRulesDisasterRecoveryConfigsClient) DisasterRecoveryConfigsListAuthorizationRulesComplete(ctx context.Context, id DisasterRecoveryConfigId) (DisasterRecoveryConfigsListAuthorizationRulesCompleteResult, error) {
	return c.DisasterRecoveryConfigsListAuthorizationRulesCompleteMatchingPredicate(ctx, id, AuthorizationRuleOperationPredicate{})
}

// DisasterRecoveryConfigsListAuthorizationRulesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AuthorizationRulesDisasterRecoveryConfigsClient) DisasterRecoveryConfigsListAuthorizationRulesCompleteMatchingPredicate(ctx context.Context, id DisasterRecoveryConfigId, predicate AuthorizationRuleOperationPredicate) (resp DisasterRecoveryConfigsListAuthorizationRulesCompleteResult, err error) {
	items := make([]AuthorizationRule, 0)

	page, err := c.DisasterRecoveryConfigsListAuthorizationRules(ctx, id)
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

	out := DisasterRecoveryConfigsListAuthorizationRulesCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForDisasterRecoveryConfigsListAuthorizationRules prepares the DisasterRecoveryConfigsListAuthorizationRules request.
func (c AuthorizationRulesDisasterRecoveryConfigsClient) preparerForDisasterRecoveryConfigsListAuthorizationRules(ctx context.Context, id DisasterRecoveryConfigId) (*http.Request, error) {
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

// preparerForDisasterRecoveryConfigsListAuthorizationRulesWithNextLink prepares the DisasterRecoveryConfigsListAuthorizationRules request with the given nextLink token.
func (c AuthorizationRulesDisasterRecoveryConfigsClient) preparerForDisasterRecoveryConfigsListAuthorizationRulesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForDisasterRecoveryConfigsListAuthorizationRules handles the response to the DisasterRecoveryConfigsListAuthorizationRules request. The method always
// closes the http.Response Body.
func (c AuthorizationRulesDisasterRecoveryConfigsClient) responderForDisasterRecoveryConfigsListAuthorizationRules(resp *http.Response) (result DisasterRecoveryConfigsListAuthorizationRulesOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result DisasterRecoveryConfigsListAuthorizationRulesOperationResponse, err error) {
			req, err := c.preparerForDisasterRecoveryConfigsListAuthorizationRulesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "authorizationrulesdisasterrecoveryconfigs.AuthorizationRulesDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsListAuthorizationRules", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "authorizationrulesdisasterrecoveryconfigs.AuthorizationRulesDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsListAuthorizationRules", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForDisasterRecoveryConfigsListAuthorizationRules(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "authorizationrulesdisasterrecoveryconfigs.AuthorizationRulesDisasterRecoveryConfigsClient", "DisasterRecoveryConfigsListAuthorizationRules", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
