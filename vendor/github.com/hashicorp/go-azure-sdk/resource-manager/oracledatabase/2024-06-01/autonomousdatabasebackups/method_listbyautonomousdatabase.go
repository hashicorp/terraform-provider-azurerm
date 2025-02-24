package autonomousdatabasebackups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByAutonomousDatabaseOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AutonomousDatabaseBackup
}

type ListByAutonomousDatabaseCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AutonomousDatabaseBackup
}

type ListByAutonomousDatabaseCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByAutonomousDatabaseCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByAutonomousDatabase ...
func (c AutonomousDatabaseBackupsClient) ListByAutonomousDatabase(ctx context.Context, id AutonomousDatabaseId) (result ListByAutonomousDatabaseOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByAutonomousDatabaseCustomPager{},
		Path:       fmt.Sprintf("%s/autonomousDatabaseBackups", id.ID()),
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
		Values *[]AutonomousDatabaseBackup `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByAutonomousDatabaseComplete retrieves all the results into a single object
func (c AutonomousDatabaseBackupsClient) ListByAutonomousDatabaseComplete(ctx context.Context, id AutonomousDatabaseId) (ListByAutonomousDatabaseCompleteResult, error) {
	return c.ListByAutonomousDatabaseCompleteMatchingPredicate(ctx, id, AutonomousDatabaseBackupOperationPredicate{})
}

// ListByAutonomousDatabaseCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AutonomousDatabaseBackupsClient) ListByAutonomousDatabaseCompleteMatchingPredicate(ctx context.Context, id AutonomousDatabaseId, predicate AutonomousDatabaseBackupOperationPredicate) (result ListByAutonomousDatabaseCompleteResult, err error) {
	items := make([]AutonomousDatabaseBackup, 0)

	resp, err := c.ListByAutonomousDatabase(ctx, id)
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

	result = ListByAutonomousDatabaseCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
