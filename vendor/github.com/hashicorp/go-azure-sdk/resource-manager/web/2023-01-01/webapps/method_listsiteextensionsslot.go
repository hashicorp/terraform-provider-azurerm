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

type ListSiteExtensionsSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SiteExtensionInfo
}

type ListSiteExtensionsSlotCompleteResult struct {
	Items []SiteExtensionInfo
}

// ListSiteExtensionsSlot ...
func (c WebAppsClient) ListSiteExtensionsSlot(ctx context.Context, id SlotId) (result ListSiteExtensionsSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/siteExtensions", id.ID()),
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
		Values *[]SiteExtensionInfo `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListSiteExtensionsSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListSiteExtensionsSlotComplete(ctx context.Context, id SlotId) (ListSiteExtensionsSlotCompleteResult, error) {
	return c.ListSiteExtensionsSlotCompleteMatchingPredicate(ctx, id, SiteExtensionInfoOperationPredicate{})
}

// ListSiteExtensionsSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListSiteExtensionsSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate SiteExtensionInfoOperationPredicate) (result ListSiteExtensionsSlotCompleteResult, err error) {
	items := make([]SiteExtensionInfo, 0)

	resp, err := c.ListSiteExtensionsSlot(ctx, id)
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

	result = ListSiteExtensionsSlotCompleteResult{
		Items: items,
	}
	return
}
