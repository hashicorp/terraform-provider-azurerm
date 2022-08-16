package backupvaults

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

type GetInSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]BackupVaultResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (GetInSubscriptionOperationResponse, error)
}

type GetInSubscriptionCompleteResult struct {
	Items []BackupVaultResource
}

func (r GetInSubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r GetInSubscriptionOperationResponse) LoadMore(ctx context.Context) (resp GetInSubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// GetInSubscription ...
func (c BackupVaultsClient) GetInSubscription(ctx context.Context, id commonids.SubscriptionId) (resp GetInSubscriptionOperationResponse, err error) {
	req, err := c.preparerForGetInSubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "GetInSubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "GetInSubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForGetInSubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "GetInSubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForGetInSubscription prepares the GetInSubscription request.
func (c BackupVaultsClient) preparerForGetInSubscription(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.DataProtection/backupVaults", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForGetInSubscriptionWithNextLink prepares the GetInSubscription request with the given nextLink token.
func (c BackupVaultsClient) preparerForGetInSubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForGetInSubscription handles the response to the GetInSubscription request. The method always
// closes the http.Response Body.
func (c BackupVaultsClient) responderForGetInSubscription(resp *http.Response) (result GetInSubscriptionOperationResponse, err error) {
	type page struct {
		Values   []BackupVaultResource `json:"value"`
		NextLink *string               `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result GetInSubscriptionOperationResponse, err error) {
			req, err := c.preparerForGetInSubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "GetInSubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "GetInSubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForGetInSubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "GetInSubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// GetInSubscriptionComplete retrieves all of the results into a single object
func (c BackupVaultsClient) GetInSubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (GetInSubscriptionCompleteResult, error) {
	return c.GetInSubscriptionCompleteMatchingPredicate(ctx, id, BackupVaultResourceOperationPredicate{})
}

// GetInSubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c BackupVaultsClient) GetInSubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate BackupVaultResourceOperationPredicate) (resp GetInSubscriptionCompleteResult, err error) {
	items := make([]BackupVaultResource, 0)

	page, err := c.GetInSubscription(ctx, id)
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

	out := GetInSubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}
