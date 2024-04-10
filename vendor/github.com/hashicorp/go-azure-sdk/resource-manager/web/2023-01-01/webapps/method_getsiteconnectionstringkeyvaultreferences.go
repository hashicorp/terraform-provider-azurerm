package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetSiteConnectionStringKeyVaultReferencesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApiKVReference
}

type GetSiteConnectionStringKeyVaultReferencesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApiKVReference
}

// GetSiteConnectionStringKeyVaultReferences ...
func (c WebAppsClient) GetSiteConnectionStringKeyVaultReferences(ctx context.Context, id commonids.AppServiceId) (result GetSiteConnectionStringKeyVaultReferencesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/config/configReferences/connectionStrings", id.ID()),
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
		Values *[]ApiKVReference `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetSiteConnectionStringKeyVaultReferencesComplete retrieves all the results into a single object
func (c WebAppsClient) GetSiteConnectionStringKeyVaultReferencesComplete(ctx context.Context, id commonids.AppServiceId) (GetSiteConnectionStringKeyVaultReferencesCompleteResult, error) {
	return c.GetSiteConnectionStringKeyVaultReferencesCompleteMatchingPredicate(ctx, id, ApiKVReferenceOperationPredicate{})
}

// GetSiteConnectionStringKeyVaultReferencesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) GetSiteConnectionStringKeyVaultReferencesCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate ApiKVReferenceOperationPredicate) (result GetSiteConnectionStringKeyVaultReferencesCompleteResult, err error) {
	items := make([]ApiKVReference, 0)

	resp, err := c.GetSiteConnectionStringKeyVaultReferences(ctx, id)
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

	result = GetSiteConnectionStringKeyVaultReferencesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
