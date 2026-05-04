package backupsv2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LongRunningBackupsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ServerBackupV2
}

type LongRunningBackupsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ServerBackupV2
}

type LongRunningBackupsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LongRunningBackupsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LongRunningBackupsList ...
func (c BackupsV2Client) LongRunningBackupsList(ctx context.Context, id FlexibleServerId) (result LongRunningBackupsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &LongRunningBackupsListCustomPager{},
		Path:       fmt.Sprintf("%s/backupsV2", id.ID()),
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
		Values *[]ServerBackupV2 `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LongRunningBackupsListComplete retrieves all the results into a single object
func (c BackupsV2Client) LongRunningBackupsListComplete(ctx context.Context, id FlexibleServerId) (LongRunningBackupsListCompleteResult, error) {
	return c.LongRunningBackupsListCompleteMatchingPredicate(ctx, id, ServerBackupV2OperationPredicate{})
}

// LongRunningBackupsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c BackupsV2Client) LongRunningBackupsListCompleteMatchingPredicate(ctx context.Context, id FlexibleServerId, predicate ServerBackupV2OperationPredicate) (result LongRunningBackupsListCompleteResult, err error) {
	items := make([]ServerBackupV2, 0)

	resp, err := c.LongRunningBackupsList(ctx, id)
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

	result = LongRunningBackupsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
