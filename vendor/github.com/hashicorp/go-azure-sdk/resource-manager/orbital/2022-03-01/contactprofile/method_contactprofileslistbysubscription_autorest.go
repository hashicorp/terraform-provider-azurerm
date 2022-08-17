package contactprofile

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactProfilesListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ContactProfile

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ContactProfilesListBySubscriptionOperationResponse, error)
}

type ContactProfilesListBySubscriptionCompleteResult struct {
	Items []ContactProfile
}

func (r ContactProfilesListBySubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ContactProfilesListBySubscriptionOperationResponse) LoadMore(ctx context.Context) (resp ContactProfilesListBySubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ContactProfilesListBySubscription ...
func (c ContactProfileClient) ContactProfilesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (resp ContactProfilesListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForContactProfilesListBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForContactProfilesListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForContactProfilesListBySubscription prepares the ContactProfilesListBySubscription request.
func (c ContactProfileClient) preparerForContactProfilesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Orbital/contactProfiles", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForContactProfilesListBySubscriptionWithNextLink prepares the ContactProfilesListBySubscription request with the given nextLink token.
func (c ContactProfileClient) preparerForContactProfilesListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForContactProfilesListBySubscription handles the response to the ContactProfilesListBySubscription request. The method always
// closes the http.Response Body.
func (c ContactProfileClient) responderForContactProfilesListBySubscription(resp *http.Response) (result ContactProfilesListBySubscriptionOperationResponse, err error) {
	type page struct {
		Values   []ContactProfile `json:"value"`
		NextLink *string          `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ContactProfilesListBySubscriptionOperationResponse, err error) {
			req, err := c.preparerForContactProfilesListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForContactProfilesListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ContactProfilesListBySubscriptionComplete retrieves all of the results into a single object
func (c ContactProfileClient) ContactProfilesListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (ContactProfilesListBySubscriptionCompleteResult, error) {
	return c.ContactProfilesListBySubscriptionCompleteMatchingPredicate(ctx, id, ContactProfileOperationPredicate{})
}

// ContactProfilesListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ContactProfileClient) ContactProfilesListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate ContactProfileOperationPredicate) (resp ContactProfilesListBySubscriptionCompleteResult, err error) {
	items := make([]ContactProfile, 0)

	page, err := c.ContactProfilesListBySubscription(ctx, id)
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

	out := ContactProfilesListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}
