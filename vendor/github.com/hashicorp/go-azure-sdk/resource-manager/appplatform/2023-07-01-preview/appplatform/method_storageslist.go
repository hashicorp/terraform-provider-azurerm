package appplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StoragesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]StorageResource
}

type StoragesListCompleteResult struct {
	Items []StorageResource
}

// StoragesList ...
func (c AppPlatformClient) StoragesList(ctx context.Context, id SpringId) (result StoragesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/storages", id.ID()),
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
		Values *[]StorageResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// StoragesListComplete retrieves all the results into a single object
func (c AppPlatformClient) StoragesListComplete(ctx context.Context, id SpringId) (StoragesListCompleteResult, error) {
	return c.StoragesListCompleteMatchingPredicate(ctx, id, StorageResourceOperationPredicate{})
}

// StoragesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) StoragesListCompleteMatchingPredicate(ctx context.Context, id SpringId, predicate StorageResourceOperationPredicate) (result StoragesListCompleteResult, err error) {
	items := make([]StorageResource, 0)

	resp, err := c.StoragesList(ctx, id)
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

	result = StoragesListCompleteResult{
		Items: items,
	}
	return
}
