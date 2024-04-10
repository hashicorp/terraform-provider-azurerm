package aad

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPolicyAssignmentListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RedisCacheAccessPolicyAssignment
}

type AccessPolicyAssignmentListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RedisCacheAccessPolicyAssignment
}

// AccessPolicyAssignmentList ...
func (c AADClient) AccessPolicyAssignmentList(ctx context.Context, id RediId) (result AccessPolicyAssignmentListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/accessPolicyAssignments", id.ID()),
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
		Values *[]RedisCacheAccessPolicyAssignment `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AccessPolicyAssignmentListComplete retrieves all the results into a single object
func (c AADClient) AccessPolicyAssignmentListComplete(ctx context.Context, id RediId) (AccessPolicyAssignmentListCompleteResult, error) {
	return c.AccessPolicyAssignmentListCompleteMatchingPredicate(ctx, id, RedisCacheAccessPolicyAssignmentOperationPredicate{})
}

// AccessPolicyAssignmentListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AADClient) AccessPolicyAssignmentListCompleteMatchingPredicate(ctx context.Context, id RediId, predicate RedisCacheAccessPolicyAssignmentOperationPredicate) (result AccessPolicyAssignmentListCompleteResult, err error) {
	items := make([]RedisCacheAccessPolicyAssignment, 0)

	resp, err := c.AccessPolicyAssignmentList(ctx, id)
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

	result = AccessPolicyAssignmentListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
