package iotdpsresource

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

type ListKeysOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SharedAccessSignatureAuthorizationRuleAccessRightsDescription

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListKeysOperationResponse, error)
}

type ListKeysCompleteResult struct {
	Items []SharedAccessSignatureAuthorizationRuleAccessRightsDescription
}

func (r ListKeysOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListKeysOperationResponse) LoadMore(ctx context.Context) (resp ListKeysOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// ListKeys ...
func (c IotDpsResourceClient) ListKeys(ctx context.Context, id commonids.ProvisioningServiceId) (resp ListKeysOperationResponse, err error) {
	req, err := c.preparerForListKeys(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListKeys", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListKeys", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListKeys(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListKeys", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListKeys prepares the ListKeys request.
func (c IotDpsResourceClient) preparerForListKeys(ctx context.Context, id commonids.ProvisioningServiceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listkeys", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListKeysWithNextLink prepares the ListKeys request with the given nextLink token.
func (c IotDpsResourceClient) preparerForListKeysWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListKeys handles the response to the ListKeys request. The method always
// closes the http.Response Body.
func (c IotDpsResourceClient) responderForListKeys(resp *http.Response) (result ListKeysOperationResponse, err error) {
	type page struct {
		Values   []SharedAccessSignatureAuthorizationRuleAccessRightsDescription `json:"value"`
		NextLink *string                                                         `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListKeysOperationResponse, err error) {
			req, err := c.preparerForListKeysWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListKeys", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListKeys", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListKeys(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "iotdpsresource.IotDpsResourceClient", "ListKeys", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListKeysComplete retrieves all of the results into a single object
func (c IotDpsResourceClient) ListKeysComplete(ctx context.Context, id commonids.ProvisioningServiceId) (ListKeysCompleteResult, error) {
	return c.ListKeysCompleteMatchingPredicate(ctx, id, SharedAccessSignatureAuthorizationRuleAccessRightsDescriptionOperationPredicate{})
}

// ListKeysCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c IotDpsResourceClient) ListKeysCompleteMatchingPredicate(ctx context.Context, id commonids.ProvisioningServiceId, predicate SharedAccessSignatureAuthorizationRuleAccessRightsDescriptionOperationPredicate) (resp ListKeysCompleteResult, err error) {
	items := make([]SharedAccessSignatureAuthorizationRuleAccessRightsDescription, 0)

	page, err := c.ListKeys(ctx, id)
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

	out := ListKeysCompleteResult{
		Items: items,
	}
	return out, nil
}
