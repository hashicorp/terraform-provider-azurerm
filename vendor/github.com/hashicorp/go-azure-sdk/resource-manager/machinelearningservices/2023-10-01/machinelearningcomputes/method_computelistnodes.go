package machinelearningcomputes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeListNodesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AmlComputeNodesInformation
}

type ComputeListNodesCompleteResult struct {
	Items []AmlComputeNodesInformation
}

// ComputeListNodes ...
func (c MachineLearningComputesClient) ComputeListNodes(ctx context.Context, id ComputeId) (result ComputeListNodesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/listNodes", id.ID()),
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
		Values *[]AmlComputeNodesInformation `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ComputeListNodesComplete retrieves all the results into a single object
func (c MachineLearningComputesClient) ComputeListNodesComplete(ctx context.Context, id ComputeId) (ComputeListNodesCompleteResult, error) {
	return c.ComputeListNodesCompleteMatchingPredicate(ctx, id, AmlComputeNodesInformationOperationPredicate{})
}

// ComputeListNodesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c MachineLearningComputesClient) ComputeListNodesCompleteMatchingPredicate(ctx context.Context, id ComputeId, predicate AmlComputeNodesInformationOperationPredicate) (result ComputeListNodesCompleteResult, err error) {
	items := make([]AmlComputeNodesInformation, 0)

	resp, err := c.ComputeListNodes(ctx, id)
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

	result = ComputeListNodesCompleteResult{
		Items: items,
	}
	return
}
