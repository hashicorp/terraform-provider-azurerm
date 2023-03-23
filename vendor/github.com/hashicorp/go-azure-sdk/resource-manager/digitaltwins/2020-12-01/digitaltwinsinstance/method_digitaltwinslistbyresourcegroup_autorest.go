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

type DigitalTwinsListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DigitalTwinsDescription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (DigitalTwinsListByResourceGroupOperationResponse, error)
}

type DigitalTwinsListByResourceGroupCompleteResult struct {
	Items []DigitalTwinsDescription
}

func (r DigitalTwinsListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r DigitalTwinsListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp DigitalTwinsListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// DigitalTwinsListByResourceGroup ...
func (c DigitalTwinsInstanceClient) DigitalTwinsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp DigitalTwinsListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForDigitalTwinsListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForDigitalTwinsListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForDigitalTwinsListByResourceGroup prepares the DigitalTwinsListByResourceGroup request.
func (c DigitalTwinsInstanceClient) preparerForDigitalTwinsListByResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
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

// preparerForDigitalTwinsListByResourceGroupWithNextLink prepares the DigitalTwinsListByResourceGroup request with the given nextLink token.
func (c DigitalTwinsInstanceClient) preparerForDigitalTwinsListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForDigitalTwinsListByResourceGroup handles the response to the DigitalTwinsListByResourceGroup request. The method always
// closes the http.Response Body.
func (c DigitalTwinsInstanceClient) responderForDigitalTwinsListByResourceGroup(resp *http.Response) (result DigitalTwinsListByResourceGroupOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result DigitalTwinsListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForDigitalTwinsListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForDigitalTwinsListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// DigitalTwinsListByResourceGroupComplete retrieves all of the results into a single object
func (c DigitalTwinsInstanceClient) DigitalTwinsListByResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (DigitalTwinsListByResourceGroupCompleteResult, error) {
	return c.DigitalTwinsListByResourceGroupCompleteMatchingPredicate(ctx, id, DigitalTwinsDescriptionOperationPredicate{})
}

// DigitalTwinsListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c DigitalTwinsInstanceClient) DigitalTwinsListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate DigitalTwinsDescriptionOperationPredicate) (resp DigitalTwinsListByResourceGroupCompleteResult, err error) {
	items := make([]DigitalTwinsDescription, 0)

	page, err := c.DigitalTwinsListByResourceGroup(ctx, id)
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

	out := DigitalTwinsListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
