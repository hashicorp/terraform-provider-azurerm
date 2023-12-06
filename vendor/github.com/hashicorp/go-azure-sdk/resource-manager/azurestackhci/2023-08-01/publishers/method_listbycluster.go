package publishers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByClusterOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Publisher
}

type ListByClusterCompleteResult struct {
	Items []Publisher
}

// ListByCluster ...
func (c PublishersClient) ListByCluster(ctx context.Context, id ClusterId) (result ListByClusterOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/publishers", id.ID()),
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
		Values *[]Publisher `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByClusterComplete retrieves all the results into a single object
func (c PublishersClient) ListByClusterComplete(ctx context.Context, id ClusterId) (ListByClusterCompleteResult, error) {
	return c.ListByClusterCompleteMatchingPredicate(ctx, id, PublisherOperationPredicate{})
}

// ListByClusterCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c PublishersClient) ListByClusterCompleteMatchingPredicate(ctx context.Context, id ClusterId, predicate PublisherOperationPredicate) (result ListByClusterCompleteResult, err error) {
	items := make([]Publisher, 0)

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
		Items: items,
	}
	return
}
