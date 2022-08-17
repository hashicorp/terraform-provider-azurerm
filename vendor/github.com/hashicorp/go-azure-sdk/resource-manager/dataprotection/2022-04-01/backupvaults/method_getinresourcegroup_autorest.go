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

type GetInResourceGroupOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]BackupVaultResource

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (GetInResourceGroupOperationResponse, error)
}

type GetInResourceGroupCompleteResult struct {
	Items []BackupVaultResource
}

func (r GetInResourceGroupOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r GetInResourceGroupOperationResponse) LoadMore(ctx context.Context) (resp GetInResourceGroupOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// GetInResourceGroup ...
func (c BackupVaultsClient) GetInResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (resp GetInResourceGroupOperationResponse, err error) {
	req, err := c.preparerForGetInResourceGroup(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "GetInResourceGroup", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "GetInResourceGroup", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForGetInResourceGroup(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "GetInResourceGroup", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForGetInResourceGroup prepares the GetInResourceGroup request.
func (c BackupVaultsClient) preparerForGetInResourceGroup(ctx context.Context, id commonids.ResourceGroupId) (*http.Request, error) {
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

// preparerForGetInResourceGroupWithNextLink prepares the GetInResourceGroup request with the given nextLink token.
func (c BackupVaultsClient) preparerForGetInResourceGroupWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForGetInResourceGroup handles the response to the GetInResourceGroup request. The method always
// closes the http.Response Body.
func (c BackupVaultsClient) responderForGetInResourceGroup(resp *http.Response) (result GetInResourceGroupOperationResponse, err error) {
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result GetInResourceGroupOperationResponse, err error) {
			req, err := c.preparerForGetInResourceGroupWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "GetInResourceGroup", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "GetInResourceGroup", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForGetInResourceGroup(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "backupvaults.BackupVaultsClient", "GetInResourceGroup", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// GetInResourceGroupComplete retrieves all of the results into a single object
func (c BackupVaultsClient) GetInResourceGroupComplete(ctx context.Context, id commonids.ResourceGroupId) (GetInResourceGroupCompleteResult, error) {
	return c.GetInResourceGroupCompleteMatchingPredicate(ctx, id, BackupVaultResourceOperationPredicate{})
}

// GetInResourceGroupCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c BackupVaultsClient) GetInResourceGroupCompleteMatchingPredicate(ctx context.Context, id commonids.ResourceGroupId, predicate BackupVaultResourceOperationPredicate) (resp GetInResourceGroupCompleteResult, err error) {
	items := make([]BackupVaultResource, 0)

	page, err := c.GetInResourceGroup(ctx, id)
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

	out := GetInResourceGroupCompleteResult{
		Items: items,
	}
	return out, nil
}
