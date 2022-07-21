package apps

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

type ListTemplatesOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]AppTemplate

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListTemplatesOperationResponse, error)
}

type ListTemplatesCompleteResult struct {
	Items []AppTemplate
}

func (r ListTemplatesOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListTemplatesOperationResponse) LoadMore(ctx context.Context) (resp ListTemplatesOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListTemplates ...
func (c AppsClient) ListTemplates(ctx context.Context, id commonids.SubscriptionId) (resp ListTemplatesOperationResponse, err error) {
	req, err := c.preparerForListTemplates(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apps.AppsClient", "ListTemplates", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "apps.AppsClient", "ListTemplates", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListTemplates(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "apps.AppsClient", "ListTemplates", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// ListTemplatesComplete retrieves all of the results into a single object
func (c AppsClient) ListTemplatesComplete(ctx context.Context, id commonids.SubscriptionId) (ListTemplatesCompleteResult, error) {
	return c.ListTemplatesCompleteMatchingPredicate(ctx, id, AppTemplateOperationPredicate{})
}

// ListTemplatesCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c AppsClient) ListTemplatesCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate AppTemplateOperationPredicate) (resp ListTemplatesCompleteResult, err error) {
	items := make([]AppTemplate, 0)

	page, err := c.ListTemplates(ctx, id)
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

	out := ListTemplatesCompleteResult{
		Items: items,
	}
	return out, nil
}

// preparerForListTemplates prepares the ListTemplates request.
func (c AppsClient) preparerForListTemplates(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.IoTCentral/appTemplates", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListTemplatesWithNextLink prepares the ListTemplates request with the given nextLink token.
func (c AppsClient) preparerForListTemplatesWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListTemplates handles the response to the ListTemplates request. The method always
// closes the http.Response Body.
func (c AppsClient) responderForListTemplates(resp *http.Response) (result ListTemplatesOperationResponse, err error) {
	type page struct {
		Values   []AppTemplate `json:"value"`
		NextLink *string       `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListTemplatesOperationResponse, err error) {
			req, err := c.preparerForListTemplatesWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "apps.AppsClient", "ListTemplates", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "apps.AppsClient", "ListTemplates", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListTemplates(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "apps.AppsClient", "ListTemplates", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}
