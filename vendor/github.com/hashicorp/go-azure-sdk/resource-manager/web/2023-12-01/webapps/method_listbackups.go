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

type ListBackupsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]BackupItem
}

type ListBackupsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []BackupItem
}

type ListBackupsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBackupsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListBackups ...
func (c WebAppsClient) ListBackups(ctx context.Context, id commonids.AppServiceId) (result ListBackupsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListBackupsCustomPager{},
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

// ListBackupsComplete retrieves all the results into a single object
func (c WebAppsClient) ListBackupsComplete(ctx context.Context, id commonids.AppServiceId) (ListBackupsCompleteResult, error) {
	return c.ListBackupsCompleteMatchingPredicate(ctx, id, BackupItemOperationPredicate{})
}

// ListBackupsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListBackupsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate BackupItemOperationPredicate) (result ListBackupsCompleteResult, err error) {
	items := make([]BackupItem, 0)

	resp, err := c.ListBackups(ctx, id)
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

	result = ListBackupsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
