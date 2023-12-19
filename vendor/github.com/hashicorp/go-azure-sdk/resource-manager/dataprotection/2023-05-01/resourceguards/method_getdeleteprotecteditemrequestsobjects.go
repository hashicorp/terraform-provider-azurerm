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

type GetDeleteProtectedItemRequestsObjectsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DppBaseResource
}

type GetDeleteProtectedItemRequestsObjectsCompleteResult struct {
	Items []DppBaseResource
}

// GetDeleteProtectedItemRequestsObjects ...
func (c ResourceGuardsClient) GetDeleteProtectedItemRequestsObjects(ctx context.Context, id ResourceGuardId) (result GetDeleteProtectedItemRequestsObjectsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/deleteProtectedItemRequests", id.ID()),
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

// GetDeleteProtectedItemRequestsObjectsComplete retrieves all the results into a single object
func (c ResourceGuardsClient) GetDeleteProtectedItemRequestsObjectsComplete(ctx context.Context, id ResourceGuardId) (GetDeleteProtectedItemRequestsObjectsCompleteResult, error) {
	return c.GetDeleteProtectedItemRequestsObjectsCompleteMatchingPredicate(ctx, id, DppBaseResourceOperationPredicate{})
}

// GetDeleteProtectedItemRequestsObjectsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceGuardsClient) GetDeleteProtectedItemRequestsObjectsCompleteMatchingPredicate(ctx context.Context, id ResourceGuardId, predicate DppBaseResourceOperationPredicate) (result GetDeleteProtectedItemRequestsObjectsCompleteResult, err error) {
	items := make([]DppBaseResource, 0)

	resp, err := c.GetDeleteProtectedItemRequestsObjects(ctx, id)
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

	result = GetDeleteProtectedItemRequestsObjectsCompleteResult{
		Items: items,
	}
	return
}
