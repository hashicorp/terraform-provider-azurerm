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

type GetUpdateProtectionPolicyRequestsObjectsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DppBaseResource
}

type GetUpdateProtectionPolicyRequestsObjectsCompleteResult struct {
	Items []DppBaseResource
}

// GetUpdateProtectionPolicyRequestsObjects ...
func (c ResourceGuardsClient) GetUpdateProtectionPolicyRequestsObjects(ctx context.Context, id ResourceGuardId) (result GetUpdateProtectionPolicyRequestsObjectsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/updateProtectionPolicyRequests", id.ID()),
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

// GetUpdateProtectionPolicyRequestsObjectsComplete retrieves all the results into a single object
func (c ResourceGuardsClient) GetUpdateProtectionPolicyRequestsObjectsComplete(ctx context.Context, id ResourceGuardId) (GetUpdateProtectionPolicyRequestsObjectsCompleteResult, error) {
	return c.GetUpdateProtectionPolicyRequestsObjectsCompleteMatchingPredicate(ctx, id, DppBaseResourceOperationPredicate{})
}

// GetUpdateProtectionPolicyRequestsObjectsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceGuardsClient) GetUpdateProtectionPolicyRequestsObjectsCompleteMatchingPredicate(ctx context.Context, id ResourceGuardId, predicate DppBaseResourceOperationPredicate) (result GetUpdateProtectionPolicyRequestsObjectsCompleteResult, err error) {
	items := make([]DppBaseResource, 0)

	resp, err := c.GetUpdateProtectionPolicyRequestsObjects(ctx, id)
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

	result = GetUpdateProtectionPolicyRequestsObjectsCompleteResult{
		Items: items,
	}
	return
}
