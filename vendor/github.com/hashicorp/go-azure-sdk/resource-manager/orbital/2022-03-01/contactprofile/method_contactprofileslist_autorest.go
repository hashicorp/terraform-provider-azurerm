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

type ContactProfilesListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]ContactProfile

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ContactProfilesListOperationResponse, error)
}

type ContactProfilesListCompleteResult struct {
	Items []ContactProfile
}

func (r ContactProfilesListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ContactProfilesListOperationResponse) LoadMore(ctx context.Context) (resp ContactProfilesListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ContactProfilesList ...
func (c ContactProfileClient) ContactProfilesList(ctx context.Context, id commonids.ResourceGroupId) (resp ContactProfilesListOperationResponse, err error) {
	req, err := c.preparerForContactProfilesList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForContactProfilesList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForContactProfilesList prepares the ContactProfilesList request.
func (c ContactProfileClient) preparerForContactProfilesList(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
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

// preparerForContactProfilesListWithNextLink prepares the ContactProfilesList request with the given nextLink token.
func (c ContactProfileClient) preparerForContactProfilesListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForContactProfilesList handles the response to the ContactProfilesList request. The method always
// closes the http.Response Body.
func (c ContactProfileClient) responderForContactProfilesList(resp *http.Response) (result ContactProfilesListOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ContactProfilesListOperationResponse, err error) {
			req, err := c.preparerForContactProfilesListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForContactProfilesList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "contactprofile.ContactProfileClient", "ContactProfilesList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ContactProfilesListComplete retrieves all of the results into a single object
func (c ContactProfileClient) ContactProfilesListComplete(ctx context.Context, id commonids.ResourceGroupId) (ContactProfilesListCompleteResult, error) {
	return c.ContactProfilesListCompleteMatchingPredicate(ctx, id, ContactProfileOperationPredicate{})
}

// ContactProfilesListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ContactProfileClient) ContactProfilesListCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate ContactProfileOperationPredicate) (resp ContactProfilesListCompleteResult, err error) {
	items := make([]ContactProfile, 0)

	page, err := c.ContactProfilesList(ctx, id)
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

	out := ContactProfilesListCompleteResult{
		Items: items,
	}
	return out, nil
}
