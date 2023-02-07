package digitaltwinsinstance

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

type DigitalTwinsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DigitalTwinsDescription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (DigitalTwinsListOperationResponse, error)
}

type DigitalTwinsListCompleteResult struct {
	Items []DigitalTwinsDescription
}

func (r DigitalTwinsListOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r DigitalTwinsListOperationResponse) LoadMore(ctx context.Context) (resp DigitalTwinsListOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// DigitalTwinsList ...
func (c DigitalTwinsInstanceClient) DigitalTwinsList(ctx context.Context, id commonids.SubscriptionId) (resp DigitalTwinsListOperationResponse, err error) {
	req, err := c.preparerForDigitalTwinsList(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsList", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsList", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForDigitalTwinsList(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsList", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForDigitalTwinsList prepares the DigitalTwinsList request.
func (c DigitalTwinsInstanceClient) preparerForDigitalTwinsList(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForDigitalTwinsListWithNextLink prepares the DigitalTwinsList request with the given nextLink token.
func (c DigitalTwinsInstanceClient) preparerForDigitalTwinsListWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForDigitalTwinsList handles the response to the DigitalTwinsList request. The method always
// closes the http.Response Body.
func (c DigitalTwinsInstanceClient) responderForDigitalTwinsList(resp *http.Response) (result DigitalTwinsListOperationResponse, err error) {
	type page struct {
		Values   []DigitalTwinsDescription `json:"value"`
		NextLink *string                   `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result DigitalTwinsListOperationResponse, err error) {
			req, err := c.preparerForDigitalTwinsListWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsList", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsList", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForDigitalTwinsList(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsList", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// DigitalTwinsListComplete retrieves all of the results into a single object
func (c DigitalTwinsInstanceClient) DigitalTwinsListComplete(ctx context.Context, id commonids.SubscriptionId) (DigitalTwinsListCompleteResult, error) {
	return c.DigitalTwinsListCompleteMatchingPredicate(ctx, id, DigitalTwinsDescriptionOperationPredicate{})
}

// DigitalTwinsListCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DigitalTwinsInstanceClient) DigitalTwinsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate DigitalTwinsDescriptionOperationPredicate) (resp DigitalTwinsListCompleteResult, err error) {
	items := make([]DigitalTwinsDescription, 0)

	page, err := c.DigitalTwinsList(ctx, id)
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

	out := DigitalTwinsListCompleteResult{
		Items: items,
	}
	return out, nil
}
