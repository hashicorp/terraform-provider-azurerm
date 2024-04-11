package applications

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

type ListByClusterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Application
}

type ListByClusterCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Application
}

// ListByCluster ...
func (c ApplicationsClient) ListByCluster(ctx context.Context, id commonids.HDInsightClusterId) (result ListByClusterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/applications", id.ID()),
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
		Values *[]Application `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByClusterComplete retrieves all the results into a single object
func (c ApplicationsClient) ListByClusterComplete(ctx context.Context, id commonids.HDInsightClusterId) (ListByClusterCompleteResult, error) {
	return c.ListByClusterCompleteMatchingPredicate(ctx, id, ApplicationOperationPredicate{})
}

// ListByClusterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApplicationsClient) ListByClusterCompleteMatchingPredicate(ctx context.Context, id commonids.HDInsightClusterId, predicate ApplicationOperationPredicate) (result ListByClusterCompleteResult, err error) {
	items := make([]Application, 0)

	resp, err := c.ListByCluster(ctx, id)
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

	result = ListByClusterCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
