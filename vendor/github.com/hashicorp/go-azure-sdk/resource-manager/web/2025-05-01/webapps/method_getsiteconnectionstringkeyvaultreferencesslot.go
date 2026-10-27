package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetSiteConnectionStringKeyVaultReferencesSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ApiKVReference
}

type GetSiteConnectionStringKeyVaultReferencesSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ApiKVReference
}

type GetSiteConnectionStringKeyVaultReferencesSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetSiteConnectionStringKeyVaultReferencesSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetSiteConnectionStringKeyVaultReferencesSlot ...
func (c WebAppsClient) GetSiteConnectionStringKeyVaultReferencesSlot(ctx context.Context, id SlotId) (result GetSiteConnectionStringKeyVaultReferencesSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetSiteConnectionStringKeyVaultReferencesSlotCustomPager{},
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

// GetSiteConnectionStringKeyVaultReferencesSlotComplete retrieves all the results into a single object
func (c WebAppsClient) GetSiteConnectionStringKeyVaultReferencesSlotComplete(ctx context.Context, id SlotId) (GetSiteConnectionStringKeyVaultReferencesSlotCompleteResult, error) {
	return c.GetSiteConnectionStringKeyVaultReferencesSlotCompleteMatchingPredicate(ctx, id, ApiKVReferenceOperationPredicate{})
}

// GetSiteConnectionStringKeyVaultReferencesSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) GetSiteConnectionStringKeyVaultReferencesSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate ApiKVReferenceOperationPredicate) (result GetSiteConnectionStringKeyVaultReferencesSlotCompleteResult, err error) {
	items := make([]ApiKVReference, 0)

	resp, err := c.GetSiteConnectionStringKeyVaultReferencesSlot(ctx, id)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
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

	result = GetSiteConnectionStringKeyVaultReferencesSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
