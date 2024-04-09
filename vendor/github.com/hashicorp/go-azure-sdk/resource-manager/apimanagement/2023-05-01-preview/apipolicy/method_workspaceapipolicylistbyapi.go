package apipolicy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiPolicyListByApiOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PolicyContract
}

type WorkspaceApiPolicyListByApiCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PolicyContract
}

// WorkspaceApiPolicyListByApi ...
func (c ApiPolicyClient) WorkspaceApiPolicyListByApi(ctx context.Context, id WorkspaceApiId) (result WorkspaceApiPolicyListByApiOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/policies", id.ID()),
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
		Values *[]PolicyContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// WorkspaceApiPolicyListByApiComplete retrieves all the results into a single object
func (c ApiPolicyClient) WorkspaceApiPolicyListByApiComplete(ctx context.Context, id WorkspaceApiId) (WorkspaceApiPolicyListByApiCompleteResult, error) {
	return c.WorkspaceApiPolicyListByApiCompleteMatchingPredicate(ctx, id, PolicyContractOperationPredicate{})
}

// WorkspaceApiPolicyListByApiCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiPolicyClient) WorkspaceApiPolicyListByApiCompleteMatchingPredicate(ctx context.Context, id WorkspaceApiId, predicate PolicyContractOperationPredicate) (result WorkspaceApiPolicyListByApiCompleteResult, err error) {
	items := make([]PolicyContract, 0)

	resp, err := c.WorkspaceApiPolicyListByApi(ctx, id)
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

	result = WorkspaceApiPolicyListByApiCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
