package resourceguards

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetUpdateProtectedItemRequestsObjectsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DppBaseResource
}

type GetUpdateProtectedItemRequestsObjectsCompleteResult struct {
	Items []DppBaseResource
}

// GetUpdateProtectedItemRequestsObjects ...
func (c ResourceGuardsClient) GetUpdateProtectedItemRequestsObjects(ctx context.Context, id ResourceGuardId) (result GetUpdateProtectedItemRequestsObjectsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/updateProtectedItemRequests", id.ID()),
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
		Values *[]DppBaseResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetUpdateProtectedItemRequestsObjectsComplete retrieves all the results into a single object
func (c ResourceGuardsClient) GetUpdateProtectedItemRequestsObjectsComplete(ctx context.Context, id ResourceGuardId) (GetUpdateProtectedItemRequestsObjectsCompleteResult, error) {
	return c.GetUpdateProtectedItemRequestsObjectsCompleteMatchingPredicate(ctx, id, DppBaseResourceOperationPredicate{})
}

// GetUpdateProtectedItemRequestsObjectsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceGuardsClient) GetUpdateProtectedItemRequestsObjectsCompleteMatchingPredicate(ctx context.Context, id ResourceGuardId, predicate DppBaseResourceOperationPredicate) (result GetUpdateProtectedItemRequestsObjectsCompleteResult, err error) {
	items := make([]DppBaseResource, 0)

	resp, err := c.GetUpdateProtectedItemRequestsObjects(ctx, id)
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

	result = GetUpdateProtectedItemRequestsObjectsCompleteResult{
		Items: items,
	}
	return
}
