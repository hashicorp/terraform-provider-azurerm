package webapps

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

type ListWorkflowsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]WorkflowEnvelope
}

type ListWorkflowsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []WorkflowEnvelope
}

// ListWorkflows ...
func (c WebAppsClient) ListWorkflows(ctx context.Context, id commonids.AppServiceId) (result ListWorkflowsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/workflows", id.ID()),
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
		Values *[]WorkflowEnvelope `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListWorkflowsComplete retrieves all the results into a single object
func (c WebAppsClient) ListWorkflowsComplete(ctx context.Context, id commonids.AppServiceId) (ListWorkflowsCompleteResult, error) {
	return c.ListWorkflowsCompleteMatchingPredicate(ctx, id, WorkflowEnvelopeOperationPredicate{})
}

// ListWorkflowsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListWorkflowsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate WorkflowEnvelopeOperationPredicate) (result ListWorkflowsCompleteResult, err error) {
	items := make([]WorkflowEnvelope, 0)

	resp, err := c.ListWorkflows(ctx, id)
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

	result = ListWorkflowsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
