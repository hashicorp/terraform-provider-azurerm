package storageinsights

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageInsightConfigsListByWorkspaceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StorageInsight
}

type StorageInsightConfigsListByWorkspaceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []StorageInsight
}

type StorageInsightConfigsListByWorkspaceCustomPager struct {
	NextLink *odata.Link `json:"@odata.nextLink"`
}

func (p *StorageInsightConfigsListByWorkspaceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// StorageInsightConfigsListByWorkspace ...
func (c StorageInsightsClient) StorageInsightConfigsListByWorkspace(ctx context.Context, id WorkspaceId) (result StorageInsightConfigsListByWorkspaceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &StorageInsightConfigsListByWorkspaceCustomPager{},
		Path:       fmt.Sprintf("%s/storageInsightConfigs", id.ID()),
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
		Values *[]StorageInsight `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// StorageInsightConfigsListByWorkspaceComplete retrieves all the results into a single object
func (c StorageInsightsClient) StorageInsightConfigsListByWorkspaceComplete(ctx context.Context, id WorkspaceId) (StorageInsightConfigsListByWorkspaceCompleteResult, error) {
	return c.StorageInsightConfigsListByWorkspaceCompleteMatchingPredicate(ctx, id, StorageInsightOperationPredicate{})
}

// StorageInsightConfigsListByWorkspaceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StorageInsightsClient) StorageInsightConfigsListByWorkspaceCompleteMatchingPredicate(ctx context.Context, id WorkspaceId, predicate StorageInsightOperationPredicate) (result StorageInsightConfigsListByWorkspaceCompleteResult, err error) {
	items := make([]StorageInsight, 0)

	resp, err := c.StorageInsightConfigsListByWorkspace(ctx, id)
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

	result = StorageInsightConfigsListByWorkspaceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
