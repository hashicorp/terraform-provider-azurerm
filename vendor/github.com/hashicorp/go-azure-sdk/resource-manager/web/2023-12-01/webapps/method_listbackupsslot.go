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

type ListBackupsSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BackupItem
}

type ListBackupsSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BackupItem
}

type ListBackupsSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBackupsSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListBackupsSlot ...
func (c WebAppsClient) ListBackupsSlot(ctx context.Context, id SlotId) (result ListBackupsSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListBackupsSlotCustomPager{},
		Path:       fmt.Sprintf("%s/backups", id.ID()),
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

// ListBackupsSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListBackupsSlotComplete(ctx context.Context, id SlotId) (ListBackupsSlotCompleteResult, error) {
	return c.ListBackupsSlotCompleteMatchingPredicate(ctx, id, BackupItemOperationPredicate{})
}

// ListBackupsSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListBackupsSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate BackupItemOperationPredicate) (result ListBackupsSlotCompleteResult, err error) {
	items := make([]BackupItem, 0)

	resp, err := c.ListBackupsSlot(ctx, id)
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

	result = ListBackupsSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
