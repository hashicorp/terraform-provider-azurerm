package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListInstanceWorkflowsSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]WorkflowEnvelope
}

type ListInstanceWorkflowsSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []WorkflowEnvelope
}

type ListInstanceWorkflowsSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListInstanceWorkflowsSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListInstanceWorkflowsSlot ...
func (c WebAppsClient) ListInstanceWorkflowsSlot(ctx context.Context, id SlotId) (result ListInstanceWorkflowsSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListInstanceWorkflowsSlotCustomPager{},
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

// ListInstanceWorkflowsSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListInstanceWorkflowsSlotComplete(ctx context.Context, id SlotId) (ListInstanceWorkflowsSlotCompleteResult, error) {
	return c.ListInstanceWorkflowsSlotCompleteMatchingPredicate(ctx, id, WorkflowEnvelopeOperationPredicate{})
}

// ListInstanceWorkflowsSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListInstanceWorkflowsSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate WorkflowEnvelopeOperationPredicate) (result ListInstanceWorkflowsSlotCompleteResult, err error) {
	items := make([]WorkflowEnvelope, 0)

	resp, err := c.ListInstanceWorkflowsSlot(ctx, id)
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

	result = ListInstanceWorkflowsSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
