package serviceresource

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServicesListByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]DataMigrationService

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ServicesListByResourceGroupOperationResponse, error)
}

type ServicesListByResourceGroupCompleteResult struct {
	Items []DataMigrationService
}

func (r ServicesListByResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ServicesListByResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp ServicesListByResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ServicesListByResourceGroup ...
func (c ServiceResourceClient) ServicesListByResourceGroup(ctx context.Context, id ResourceGroupId) (resp ServicesListByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForServicesListByResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesListByResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForServicesListByResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesListByResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForServicesListByResourceGroup prepares the ServicesListByResourceGroup request.
func (c ServiceResourceClient) preparerForServicesListByResourceGroup(ctx context.Context, id ResourceGroupId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.DataMigration/services", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForServicesListByResourceGroupWithNextLink prepares the ServicesListByResourceGroup request with the given nextLink token.
func (c ServiceResourceClient) preparerForServicesListByResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForServicesListByResourceGroup handles the response to the ServicesListByResourceGroup request. The method always
// closes the http.Response Body.
func (c ServiceResourceClient) responderForServicesListByResourceGroup(resp *http.Response) (result ServicesListByResourceGroupOperationResponse, err error) {
	type page struct {
		Values   []DataMigrationService `json:"value"`
		NextLink *string                `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ServicesListByResourceGroupOperationResponse, err error) {
			req, err := c.preparerForServicesListByResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesListByResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesListByResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForServicesListByResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "serviceresource.ServiceResourceClient", "ServicesListByResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ServicesListByResourceGroupComplete retrieves all of the results into a single object
func (c ServiceResourceClient) ServicesListByResourceGroupComplete(ctx context.Context, id ResourceGroupId) (ServicesListByResourceGroupCompleteResult, error) {
	return c.ServicesListByResourceGroupCompleteMatchingPredicate(ctx, id, DataMigrationServiceOperationPredicate{})
}

// ServicesListByResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c ServiceResourceClient) ServicesListByResourceGroupCompleteMatchingPredicate(ctx context.Context, id ResourceGroupId, predicate DataMigrationServiceOperationPredicate) (resp ServicesListByResourceGroupCompleteResult, err error) {
	items := make([]DataMigrationService, 0)

	page, err := c.ServicesListByResourceGroup(ctx, id)
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

	out := ServicesListByResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
