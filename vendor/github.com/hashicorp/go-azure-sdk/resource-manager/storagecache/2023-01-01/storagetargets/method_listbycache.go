package storagetargets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByCacheOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StorageTarget
}

type ListByCacheCompleteResult struct {
	Items []StorageTarget
}

// ListByCache ...
func (c StorageTargetsClient) ListByCache(ctx context.Context, id CacheId) (result ListByCacheOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/storageTargets", id.ID()),
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
		Values *[]StorageTarget `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByCacheComplete retrieves all the results into a single object
func (c StorageTargetsClient) ListByCacheComplete(ctx context.Context, id CacheId) (ListByCacheCompleteResult, error) {
	return c.ListByCacheCompleteMatchingPredicate(ctx, id, StorageTargetOperationPredicate{})
}

// ListByCacheCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StorageTargetsClient) ListByCacheCompleteMatchingPredicate(ctx context.Context, id CacheId, predicate StorageTargetOperationPredicate) (result ListByCacheCompleteResult, err error) {
	items := make([]StorageTarget, 0)

	resp, err := c.ListByCache(ctx, id)
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

	result = ListByCacheCompleteResult{
		Items: items,
	}
	return
}
