package managedidentity

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type UserAssignedIdentitiesListBySubscriptionResponse struct {
	HttpResponse *http.Response
	Model        *[]Identity

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (UserAssignedIdentitiesListBySubscriptionResponse, error)
}

type UserAssignedIdentitiesListBySubscriptionCompleteResult struct {
	Items []Identity
}

func (r UserAssignedIdentitiesListBySubscriptionResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r UserAssignedIdentitiesListBySubscriptionResponse) LoadMore(ctx context.Context) (resp UserAssignedIdentitiesListBySubscriptionResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// UserAssignedIdentitiesListBySubscription ...
func (c ManagedIdentityClient) UserAssignedIdentitiesListBySubscription(ctx context.Context, id SubscriptionId) (resp UserAssignedIdentitiesListBySubscriptionResponse, err error) {
	req, err := c.preparerForUserAssignedIdentitiesListBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForUserAssignedIdentitiesListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// UserAssignedIdentitiesListBySubscriptionComplete retrieves all of the results into a single object
func (c ManagedIdentityClient) UserAssignedIdentitiesListBySubscriptionComplete(ctx context.Context, id SubscriptionId) (UserAssignedIdentitiesListBySubscriptionCompleteResult, error) {
	return c.UserAssignedIdentitiesListBySubscriptionCompleteMatchingPredicate(ctx, id, IdentityPredicate{})
}

// UserAssignedIdentitiesListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ManagedIdentityClient) UserAssignedIdentitiesListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id SubscriptionId, predicate IdentityPredicate) (resp UserAssignedIdentitiesListBySubscriptionCompleteResult, err error) {
	items := make([]Identity, 0)

	page, err := c.UserAssignedIdentitiesListBySubscription(ctx, id)
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

	out := UserAssignedIdentitiesListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForUserAssignedIdentitiesListBySubscription prepares the UserAssignedIdentitiesListBySubscription request.
func (c ManagedIdentityClient) preparerForUserAssignedIdentitiesListBySubscription(ctx context.Context, id SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForUserAssignedIdentitiesListBySubscriptionWithNextLink prepares the UserAssignedIdentitiesListBySubscription request with the given nextLink token.
func (c ManagedIdentityClient) preparerForUserAssignedIdentitiesListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForUserAssignedIdentitiesListBySubscription handles the response to the UserAssignedIdentitiesListBySubscription request. The method always
// closes the http.Response Body.
func (c ManagedIdentityClient) responderForUserAssignedIdentitiesListBySubscription(resp *http.Response) (result UserAssignedIdentitiesListBySubscriptionResponse, err error) {
	type page struct {
		Values   []Identity `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result UserAssignedIdentitiesListBySubscriptionResponse, err error) {
			req, err := c.preparerForUserAssignedIdentitiesListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForUserAssignedIdentitiesListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "managedidentity.ManagedIdentityClient", "UserAssignedIdentitiesListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
