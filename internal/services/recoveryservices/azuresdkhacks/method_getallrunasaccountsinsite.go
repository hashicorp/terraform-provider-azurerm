// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azuresdkhacks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/runasaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// workaround for https://github.com/Azure/azure-rest-api-specs/issues/24712
// the difference is in the struct `RunAsAccountProperties`
// TODO 4.0: check if this could be removed

type RunAsAccountsClient struct {
	Client *resourcemanager.Client
}

type GetAllRunAsAccountsInSiteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VMwareRunAsAccount
}

type GetAllRunAsAccountsInSiteCompleteResult struct {
	Items []VMwareRunAsAccount
}

// GetAllRunAsAccountsInSite ...
func (c RunAsAccountsClient) GetAllRunAsAccountsInSite(ctx context.Context, id runasaccounts.VMwareSiteId) (result GetAllRunAsAccountsInSiteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json",
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
func (c RunAsAccountsClient) GetAllRunAsAccountsInSiteComplete(ctx context.Context, id runasaccounts.VMwareSiteId) (GetAllRunAsAccountsInSiteCompleteResult, error) {
	items := make([]VMwareRunAsAccount, 0)

	resp, err := c.GetAllRunAsAccountsInSite(ctx, id)
	if err != nil {
		err = fmt.Errorf("loading results: %+v", err)
		return GetAllRunAsAccountsInSiteCompleteResult{}, err
	}
	if resp.Model != nil {
		items = append(items, *resp.Model...)
	}

	return GetAllRunAsAccountsInSiteCompleteResult{Items: items}, nil
}
