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

type ListSiteBackupsSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BackupItem
}

type ListSiteBackupsSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BackupItem
}

// ListSiteBackupsSlot ...
func (c WebAppsClient) ListSiteBackupsSlot(ctx context.Context, id SlotId) (result ListSiteBackupsSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/listbackups", id.ID()),
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
		Values *[]BackupItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListSiteBackupsSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListSiteBackupsSlotComplete(ctx context.Context, id SlotId) (ListSiteBackupsSlotCompleteResult, error) {
	return c.ListSiteBackupsSlotCompleteMatchingPredicate(ctx, id, BackupItemOperationPredicate{})
}

// ListSiteBackupsSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListSiteBackupsSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate BackupItemOperationPredicate) (result ListSiteBackupsSlotCompleteResult, err error) {
	items := make([]BackupItem, 0)

	resp, err := c.ListSiteBackupsSlot(ctx, id)
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

	result = ListSiteBackupsSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
