package nodetype

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByManagedClustersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NodeType
}

type ListByManagedClustersCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []NodeType
}

type ListByManagedClustersCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByManagedClustersCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByManagedClusters ...
func (c NodeTypeClient) ListByManagedClusters(ctx context.Context, id ManagedClusterId) (result ListByManagedClustersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByManagedClustersCustomPager{},
		Path:       fmt.Sprintf("%s/nodeTypes", id.ID()),
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
		Values *[]NodeType `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByManagedClustersComplete retrieves all the results into a single object
func (c NodeTypeClient) ListByManagedClustersComplete(ctx context.Context, id ManagedClusterId) (ListByManagedClustersCompleteResult, error) {
	return c.ListByManagedClustersCompleteMatchingPredicate(ctx, id, NodeTypeOperationPredicate{})
}

// ListByManagedClustersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NodeTypeClient) ListByManagedClustersCompleteMatchingPredicate(ctx context.Context, id ManagedClusterId, predicate NodeTypeOperationPredicate) (result ListByManagedClustersCompleteResult, err error) {
	items := make([]NodeType, 0)

	resp, err := c.ListByManagedClusters(ctx, id)
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

	result = ListByManagedClustersCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
