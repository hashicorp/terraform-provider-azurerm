package runasaccounts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetAllRunAsAccountsInSiteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VMwareRunAsAccount
}

type GetAllRunAsAccountsInSiteCompleteResult struct {
	Items []VMwareRunAsAccount
}

// GetAllRunAsAccountsInSite ...
func (c RunAsAccountsClient) GetAllRunAsAccountsInSite(ctx context.Context, id VMwareSiteId) (result GetAllRunAsAccountsInSiteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/runAsAccounts", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]VMwareRunAsAccount `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetAllRunAsAccountsInSiteComplete retrieves all the results into a single object
func (c RunAsAccountsClient) GetAllRunAsAccountsInSiteComplete(ctx context.Context, id VMwareSiteId) (GetAllRunAsAccountsInSiteCompleteResult, error) {
	return c.GetAllRunAsAccountsInSiteCompleteMatchingPredicate(ctx, id, VMwareRunAsAccountOperationPredicate{})
}

// GetAllRunAsAccountsInSiteCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c RunAsAccountsClient) GetAllRunAsAccountsInSiteCompleteMatchingPredicate(ctx context.Context, id VMwareSiteId, predicate VMwareRunAsAccountOperationPredicate) (result GetAllRunAsAccountsInSiteCompleteResult, err error) {
	items := make([]VMwareRunAsAccount, 0)

	resp, err := c.GetAllRunAsAccountsInSite(ctx, id)
	if err != nil {
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = GetAllRunAsAccountsInSiteCompleteResult{
		Items: items,
	}
	return
}
