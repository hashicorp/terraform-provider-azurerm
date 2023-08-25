package labs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListVhdsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LabVhd
}

type ListVhdsCompleteResult struct {
	Items []LabVhd
}

// ListVhds ...
func (c LabsClient) ListVhds(ctx context.Context, id LabId) (result ListVhdsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Path:       fmt.Sprintf("%s/listVhds", id.ID()),
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
		Values *[]LabVhd `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListVhdsComplete retrieves all the results into a single object
func (c LabsClient) ListVhdsComplete(ctx context.Context, id LabId) (ListVhdsCompleteResult, error) {
	return c.ListVhdsCompleteMatchingPredicate(ctx, id, LabVhdOperationPredicate{})
}

// ListVhdsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LabsClient) ListVhdsCompleteMatchingPredicate(ctx context.Context, id LabId, predicate LabVhdOperationPredicate) (result ListVhdsCompleteResult, err error) {
	items := make([]LabVhd, 0)

	resp, err := c.ListVhds(ctx, id)
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

	result = ListVhdsCompleteResult{
		Items: items,
	}
	return
}
