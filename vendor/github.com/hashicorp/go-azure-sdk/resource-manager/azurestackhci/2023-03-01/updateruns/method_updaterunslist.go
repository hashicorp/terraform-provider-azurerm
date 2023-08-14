package updateruns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateRunsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]UpdateRun
}

type UpdateRunsListCompleteResult struct {
	Items []UpdateRun
}

// UpdateRunsList ...
func (c UpdateRunsClient) UpdateRunsList(ctx context.Context, id UpdateId) (result UpdateRunsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/updateRuns", id.ID()),
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
		Values *[]UpdateRun `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// UpdateRunsListComplete retrieves all the results into a single object
func (c UpdateRunsClient) UpdateRunsListComplete(ctx context.Context, id UpdateId) (UpdateRunsListCompleteResult, error) {
	return c.UpdateRunsListCompleteMatchingPredicate(ctx, id, UpdateRunOperationPredicate{})
}

// UpdateRunsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c UpdateRunsClient) UpdateRunsListCompleteMatchingPredicate(ctx context.Context, id UpdateId, predicate UpdateRunOperationPredicate) (result UpdateRunsListCompleteResult, err error) {
	items := make([]UpdateRun, 0)

	resp, err := c.UpdateRunsList(ctx, id)
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

	result = UpdateRunsListCompleteResult{
		Items: items,
	}
	return
}
